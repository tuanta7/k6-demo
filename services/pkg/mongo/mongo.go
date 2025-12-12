package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	URI            string
	Database       string
	MaxPoolSize    uint64
	MinPoolSize    uint64
	ConnectTimeout time.Duration
	QueryTimeout   time.Duration
}

type Client struct {
	client   *mongo.Client
	database *mongo.Database
	timeout  time.Duration
}

func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	opts := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetConnectTimeout(cfg.ConnectTimeout)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	if pingErr := client.Ping(ctx, readpref.Primary()); pingErr != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", pingErr)
	}

	return &Client{
		client:   client,
		database: client.Database(cfg.Database),
		timeout:  cfg.QueryTimeout,
	}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c *Client) Get(ctx context.Context, collection string, filter any, result any) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.database.Collection(collection).FindOne(ctx, filter).Decode(result)
}
