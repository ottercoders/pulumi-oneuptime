package resources

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ottercoders/pulumi-oneuptime/provider/client"
	"github.com/pulumi/pulumi-go-provider/infer"
)

type Config struct {
	ApiKey    string  `pulumi:"apiKey,optional"`
	BaseURL   *string `pulumi:"baseUrl,optional"`
	ProjectID *string `pulumi:"projectId,optional"`

	client *client.Client
}

var _ infer.Annotated = (*Config)(nil)
var _ infer.CustomConfigure = (*Config)(nil)

func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.ApiKey, "The API key for authenticating with OneUptime. Can also be set via ONEUPTIME_API_KEY env var.")
	a.Describe(&c.BaseURL, "The base URL of the OneUptime instance. Defaults to https://oneuptime.com. Can also be set via ONEUPTIME_BASE_URL env var.")
	a.Describe(&c.ProjectID, "The default project ID for resources. Can also be set via ONEUPTIME_PROJECT_ID env var.")
}

func (c *Config) Configure(ctx context.Context) error {
	if c.ApiKey == "" {
		c.ApiKey = os.Getenv("ONEUPTIME_API_KEY")
	}
	if c.ApiKey == "" {
		return infer.ProviderErrorf("apiKey is required: set it in provider config or ONEUPTIME_API_KEY env var")
	}

	baseURL := "https://oneuptime.com"
	if c.BaseURL != nil && *c.BaseURL != "" {
		baseURL = strings.TrimRight(*c.BaseURL, "/")
	} else if env := os.Getenv("ONEUPTIME_BASE_URL"); env != "" {
		baseURL = strings.TrimRight(env, "/")
	}

	if c.ProjectID == nil {
		if env := os.Getenv("ONEUPTIME_PROJECT_ID"); env != "" {
			c.ProjectID = &env
		}
	}

	c.client = &client.Client{
		BaseURL: baseURL,
		APIKey:  c.ApiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return nil
}

func (c *Config) GetClient() *client.Client {
	return c.client
}
