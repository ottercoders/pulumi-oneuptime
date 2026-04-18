package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sync"
)

type Client struct {
	BaseURL    string
	APIKey     string
	TenantID   string
	HTTPClient *http.Client

	// unknownCols tracks columns that the server has rejected from a select
	// clause (HTTP 400 "TableColumnMetadata not found for X column"). Populated
	// lazily so subsequent reads against the same endpoint skip the failing
	// column upfront instead of paying a round trip to discover it again.
	colMu       sync.RWMutex
	unknownCols map[string]map[string]bool
}

type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("oneuptime API error (HTTP %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("oneuptime API error (HTTP %d): %s", e.StatusCode, e.Body)
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == 404
}

// unknownColumnRegex matches OneUptime's schema-mismatch error message.
var unknownColumnRegex = regexp.MustCompile(`TableColumnMetadata not found for (\w+) column`)

func extractUnknownColumn(msg string) string {
	m := unknownColumnRegex.FindStringSubmatch(msg)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

func (c *Client) markUnknownColumn(path, col string) {
	c.colMu.Lock()
	defer c.colMu.Unlock()
	if c.unknownCols == nil {
		c.unknownCols = map[string]map[string]bool{}
	}
	if c.unknownCols[path] == nil {
		c.unknownCols[path] = map[string]bool{}
	}
	c.unknownCols[path][col] = true
}

// stripUnknownColumns returns a copy of sel with any columns that have
// previously been rejected for this path removed.
func (c *Client) stripUnknownColumns(path string, sel map[string]bool) map[string]bool {
	if sel == nil {
		return nil
	}
	c.colMu.RLock()
	blocked := c.unknownCols[path]
	c.colMu.RUnlock()
	out := make(map[string]bool, len(sel))
	for k, v := range sel {
		if !blocked[k] {
			out[k] = v
		}
	}
	return out
}

// stripUnknownColumnsFromData returns a copy of data with any columns that
// have previously been rejected for this path removed. Used on Create/Update
// bodies so writes don't carry fields the target OneUptime build refuses.
func (c *Client) stripUnknownColumnsFromData(path string, data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}
	c.colMu.RLock()
	blocked := c.unknownCols[path]
	c.colMu.RUnlock()
	out := make(map[string]interface{}, len(data))
	for k, v := range data {
		if !blocked[k] {
			out[k] = v
		}
	}
	return out
}

// selectWithTolerance invokes doCall with a select map, and on a 400 response
// naming an unknown column, drops that column and retries. Lets reads succeed
// against self-hosted OneUptime builds whose schema is behind the upstream set
// of columns the provider's struct tags request.
func (c *Client) selectWithTolerance(path string, sel map[string]bool, doCall func(map[string]bool) (map[string]interface{}, error)) (map[string]interface{}, error) {
	current := c.stripUnknownColumns(path, sel)
	const maxAttempts = 8
	for i := 0; i < maxAttempts; i++ {
		resp, err := doCall(current)
		if err == nil {
			return resp, nil
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) || apiErr.StatusCode != 400 {
			return nil, err
		}
		col := extractUnknownColumn(apiErr.Message)
		if col == "" || current == nil || !current[col] {
			return nil, err
		}
		c.markUnknownColumn(path, col)
		delete(current, col)
	}
	return nil, fmt.Errorf("%s: too many unknown-column retries", path)
}

// dataWithTolerance is the Create/Update counterpart to selectWithTolerance:
// on a 400 "TableColumnMetadata not found" from the server, drop that column
// from the body and retry. Silently discards user-supplied values the server
// doesn't recognize — surfaces a debug log when ONEUPTIME_DEBUG is set so
// callers can see what was dropped.
func (c *Client) dataWithTolerance(path string, data map[string]interface{}, doCall func(map[string]interface{}) (map[string]interface{}, error)) (map[string]interface{}, error) {
	current := c.stripUnknownColumnsFromData(path, data)
	const maxAttempts = 8
	for i := 0; i < maxAttempts; i++ {
		resp, err := doCall(current)
		if err == nil {
			return resp, nil
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) || apiErr.StatusCode != 400 {
			return nil, err
		}
		col := extractUnknownColumn(apiErr.Message)
		if col == "" || current == nil {
			return nil, err
		}
		if _, present := current[col]; !present {
			return nil, err
		}
		if os.Getenv("ONEUPTIME_DEBUG") != "" {
			fmt.Fprintf(os.Stderr, "[oneuptime] %s: dropping unknown column %q from write body and retrying\n", path, col)
		}
		c.markUnknownColumn(path, col)
		delete(current, col)
	}
	return nil, fmt.Errorf("%s: too many unknown-column retries", path)
}

// tenantFromData returns a tenantid override to use for a request whose body
// carries a projectId. OneUptime's create/update endpoints enforce scoping via
// the tenantid header; when the provider config omits projectId but the
// resource arg supplies it, the header would otherwise be empty.
func tenantFromData(data map[string]interface{}) string {
	if s, ok := data["projectId"].(string); ok {
		return s
	}
	return ""
}

func (c *Client) CreateResource(ctx context.Context, path string, data map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/%s", c.BaseURL, path)

	return c.dataWithTolerance(path, data, func(d map[string]interface{}) (map[string]interface{}, error) {
		// OneUptime API expects create bodies wrapped in {"data": {...}}
		body, err := json.Marshal(map[string]interface{}{"data": d})
		if err != nil {
			return nil, fmt.Errorf("marshalling create body: %w", err)
		}
		return c.doRequest(ctx, http.MethodPost, url, body, tenantFromData(d))
	})
}

func (c *Client) ReadResource(ctx context.Context, path string, id string, selectFields map[string]bool) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/%s/%s/get-item", c.BaseURL, path, id)

	return c.selectWithTolerance(path, selectFields, func(sel map[string]bool) (map[string]interface{}, error) {
		reqBody := map[string]interface{}{}
		if sel != nil {
			reqBody["select"] = sel
		}
		body, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("marshalling read body: %w", err)
		}
		return c.doRequest(ctx, http.MethodPost, url, body, "")
	})
}

func (c *Client) UpdateResource(ctx context.Context, path string, id string, data map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/%s/%s", c.BaseURL, path, id)

	_, err := c.dataWithTolerance(path, data, func(d map[string]interface{}) (map[string]interface{}, error) {
		// OneUptime API expects update bodies wrapped in {"data": {...}}
		body, err := json.Marshal(map[string]interface{}{"data": d})
		if err != nil {
			return nil, fmt.Errorf("marshalling update body: %w", err)
		}
		return c.doRequest(ctx, http.MethodPut, url, body, tenantFromData(d))
	})
	return err
}

func (c *Client) ListResources(ctx context.Context, path string, query map[string]interface{}, selectFields map[string]bool, limit int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/%s/get-list", c.BaseURL, path)

	resp, err := c.selectWithTolerance(path, selectFields, func(sel map[string]bool) (map[string]interface{}, error) {
		reqBody := map[string]interface{}{
			"limit": limit,
			"skip":  0,
		}
		if query != nil {
			reqBody["query"] = query
		}
		if sel != nil {
			reqBody["select"] = sel
		}
		body, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("marshalling list body: %w", err)
		}
		return c.doRequest(ctx, http.MethodPost, url, body, "")
	})
	if err != nil {
		return nil, err
	}

	dataRaw, ok := resp["data"]
	if !ok {
		return nil, nil
	}

	dataArr, ok := dataRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected data array in list response, got %T", dataRaw)
	}

	results := make([]map[string]interface{}, 0, len(dataArr))
	for _, item := range dataArr {
		if m, ok := item.(map[string]interface{}); ok {
			results = append(results, m)
		}
	}

	return results, nil
}

func (c *Client) DeleteResource(ctx context.Context, path string, id string) error {
	url := fmt.Sprintf("%s/api/%s/%s", c.BaseURL, path, id)

	_, err := c.doRequest(ctx, http.MethodDelete, url, nil, "")
	if err != nil {
		if IsNotFound(err) {
			return nil // Already deleted
		}
		return err
	}
	return nil
}

// doRequest issues an HTTP request. If tenantOverride is non-empty it is sent
// as the tenantid header instead of c.TenantID, so create/update endpoints can
// scope to a projectId supplied on the resource arg rather than only on the
// provider config.
func (c *Client) doRequest(ctx context.Context, method, url string, body []byte, tenantOverride string) (map[string]interface{}, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("APIKey", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	tenantID := c.TenantID
	if tenantOverride != "" {
		tenantID = tenantOverride
	}
	if tenantID != "" {
		req.Header.Set("tenantid", tenantID)
	}

	debug := os.Getenv("ONEUPTIME_DEBUG") != ""
	if debug {
		fmt.Fprintf(os.Stderr, "[oneuptime] %s %s body=%s\n", method, url, string(body))
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", method, url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if debug {
		fmt.Fprintf(os.Stderr, "[oneuptime] response status=%d body=%s\n", resp.StatusCode, string(respBody))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
		// Try to extract error message from JSON response
		var errResp map[string]interface{}
		if json.Unmarshal(respBody, &errResp) == nil {
			if msg, ok := errResp["message"].(string); ok {
				apiErr.Message = msg
			} else if msg, ok := errResp["error"].(string); ok {
				apiErr.Message = msg
			}
		}
		return nil, apiErr
	}

	// Some endpoints return empty body (e.g., DELETE, PUT)
	// Also treat "{}" as empty (e.g., Read after delete returns 200 with {})
	if len(respBody) == 0 || string(respBody) == "{}" {
		return nil, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return result, nil
}
