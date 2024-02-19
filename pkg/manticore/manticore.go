package manticore

import (
	"context"
	"fmt"

	"github.com/glu-project/configs"

	"github.com/manticoresoftware/go-sdk/manticore"
)

type ManticoreClient struct {
	*manticore.Client
	Config *configs.Database
}

func NewClient(config *configs.Database) (*ManticoreClient, error) {
	return &ManticoreClient{
		Config: config,
	}, nil
}

func (c *ManticoreClient) Connect(ctx context.Context) error {
	cl := manticore.NewClient()
	cl.SetServer(c.Config.Host, uint16(c.Config.Port))
	c.Client = &cl

	if _, err := c.Client.Open(); err != nil {
		return fmt.Errorf("unable to connect manticore: %w", err)
	}

	if _, err := c.Client.Ping(123); err != nil {
		return fmt.Errorf("unable to connect manticore: %w", err)
	}

	return nil
}

func (c *ManticoreClient) Close(ctx context.Context) error {
	if _, err := c.Client.Close(); err != nil {
		return fmt.Errorf("unable to close Manticore: %w", err)
	}

	return nil
}
