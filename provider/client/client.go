package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
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

func (c *Client) CreateResource(ctx context.Context, path string, data map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/%s", c.BaseURL, path)

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshalling create body: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) ReadResource(ctx context.Context, path string, id string, selectFields map[string]bool) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/%s/%s/get-item", c.BaseURL, path, id)

	reqBody := map[string]interface{}{}
	if selectFields != nil {
		reqBody["select"] = selectFields
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling read body: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) UpdateResource(ctx context.Context, path string, id string, data map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/%s/%s", c.BaseURL, path, id)

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshalling update body: %w", err)
	}

	_, err = c.doRequest(ctx, http.MethodPut, url, body)
	return err
}

func (c *Client) DeleteResource(ctx context.Context, path string, id string) error {
	url := fmt.Sprintf("%s/api/%s/%s", c.BaseURL, path, id)

	_, err := c.doRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		if IsNotFound(err) {
			return nil // Already deleted
		}
		return err
	}
	return nil
}

func (c *Client) doRequest(ctx context.Context, method, url string, body []byte) (map[string]interface{}, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("ApiKey", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", method, url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
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
	if len(respBody) == 0 {
		return nil, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return result, nil
}
