package http_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/glu-project/configs"
	"io"
	"log"
	"net/http"
	"net/url"
)

const CHUNK_SIZE = 1024

type HttpClient struct {
	*http.Client
	Endpoint configs.Endpoint
}

func NewHttpClient(endpoint configs.Endpoint) *HttpClient {
	return &HttpClient{
		Endpoint: endpoint,
		Client:   http.DefaultClient,
	}
}

func UnmarshalResponse[T any](data []byte) (*T, error) {
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response: %w", err)
	}

	return &result, nil
}

func (c *HttpClient) GET(ctx context.Context, paths ...string) ([]byte, error) {
	targetUrl, err := url.JoinPath(c.Endpoint.Host, paths...)
	if err != nil {
		return nil, fmt.Errorf("unable to join path: %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	if c.Endpoint.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Endpoint.Token))
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to call request api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK && resp.StatusCode > http.StatusIMUsed {
		return nil, fmt.Errorf("failed to retrieve data with status %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %w", err)
	}

	return body, nil
}

func (c *HttpClient) GETWithIo(ctx context.Context, paths ...string) (io.ReadCloser, error) {
	targetUrl, err := url.JoinPath(c.Endpoint.Host, paths...)
	if err != nil {
		return nil, fmt.Errorf("unable to join path: %w", err)
	}
	log.Println("targetUrl", targetUrl)
	req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	if c.Endpoint.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Endpoint.Token))
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to call request api: %w", err)
	}

	return resp.Body, nil
}
