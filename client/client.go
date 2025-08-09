package client

import (
	"context"
	"fmt"

	xai "github.com/Maniacal/go-xai-sdk/xai/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn   *grpc.ClientConn
	Models xai.ModelsClient
}

// New creates a new Client.
func New(apiKey, endpoint string) (*Client, error) {
	// Dial gRPC with TLS.
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), grpc.WithUnaryInterceptor(authInterceptor(apiKey)))
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	return &Client{
		conn:   conn,
		Models: xai.NewModelsClient(conn),
	}, nil
}

// authInterceptor adds API key to metadata.
func authInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", apiKey)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// Close the connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
