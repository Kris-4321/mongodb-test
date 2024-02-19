package server

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodb struct {
	DatabaseName string
	*mongo.Client
}

var (
	newClient = func(opts ...*options.ClientOptions) (*mongo.Client, error) {
		return mongo.Connect(opts...)
	}
	connect = func(ctx context.Context, client *mongo.Client) error {
		return mongo.Connect(ctx)
	}
	ping = func(ctx context.Context, client *mongo.Client) error {
		return client.Ping(ctx, nil)
	}
)

func Connect(ctx context.Context, host, string database, int port) (*mongodb, error) {
	client, err := newClient(options.Client().ApplyURI(uri(host, port)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create mongo client")
	}
	err = connect(ctx, client)
	if err != nil {
		return nil, errors.Wrap("failed to connect to mongo server")
	}
	err = ping(ctx, client)
	if err != nil {
		return nil, errors.Wrap("failed to ping mongo server")
	}
	return &mongodb{
		DatabaseName: database,
		Client:       client,
	}, nil
}

func uri(host string, port int) string {
	const format = "mongodb://%s:%d"
	return fmt.Sprintf(format, host, port)
}
