package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is a base client for gRPC services
type Client struct {
	conn *grpc.ClientConn
}

// NewClient creates a new gRPC client with a connection to the specified address
func NewClient(address string) (*Client, error) {
	// Create a connection without blocking
	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Connection returns the underlying gRPC connection
func (c *Client) Connection() *grpc.ClientConn {
	return c.conn
}