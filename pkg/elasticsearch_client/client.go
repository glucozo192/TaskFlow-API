package elastic_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/glu-project/configs"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticClient struct {
	*elasticsearch.Client
	ElaCfg elasticsearch.Config
	Cfg    *configs.Database
}

func NewElasticClient(cfg *configs.Database) (*ElasticClient, error) {
	elaCfg := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s", cfg.Address())},
		Username:  cfg.Username,
		Password:  cfg.Password,
	}
	return &ElasticClient{
		Cfg:    cfg,
		ElaCfg: elaCfg,
	}, nil
}

func (c *ElasticClient) Connect(ctx context.Context) error {
	cfg := elasticsearch.Config{
		// using fasthttp transport
		Transport: new(transport),
		Addresses: []string{fmt.Sprintf("http://%s", c.Cfg.Address())},
		Username:  c.Cfg.Username,
		Password:  c.Cfg.Password,
		// Transport: &http.Transport{
		// 	MaxIdleConnsPerHost:   10,
		// 	ResponseHeaderTimeout: time.Second,
		// 	DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
		// 	TLSClientConfig: &tls.Config{
		// 		MinVersion: tls.VersionTLS12,
		// 	},
		// },
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("connect elastic error: %w", err)
	}
	if _, err := client.Ping(); err != nil {
		return fmt.Errorf("ping elastic error: %w", err)
	}

	c.Client = client
	return nil
}

func (c *ElasticClient) Close(ctx context.Context) error {
	return nil
}

func (c *ElasticClient) CreateDocument(ctx context.Context, index string, content any) error {
	query := map[string]any{
		"query": map[string]any{
			"match": content,
		},
	}
	b, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("unable to marshal content: %w", err)
	}
	req := esapi.IndexRequest{
		Index:   index,
		Body:    bytes.NewReader(b),
		Refresh: "true",
		Pretty:  true,
	}
	res, err := req.Do(ctx, c.Client)
	if err != nil {
		return fmt.Errorf("unable to getting response: %w", err)
	}

	defer res.Body.Close()
	return nil
}
