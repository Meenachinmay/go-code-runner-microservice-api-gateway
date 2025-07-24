package executor

import (
	"context"
	"fmt"
	executorpb "go-code-runner-microservice/api-gateway/go-code-runner-microservice/proto/executor/v1"
	baseClient "go-code-runner-microservice/api-gateway/internal/service/grpc"
	"google.golang.org/grpc"
)

type Client struct {
	client executorpb.ExecutorServiceClient
	base   *baseClient.Client
}

func NewClient(address string) (*Client, error) {
	base, err := baseClient.NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor client: %w", err)
	}

	return &Client{
		client: executorpb.NewExecutorServiceClient(base.Connection()),
		base:   base,
	}, nil
}

func NewClientWithOptions(address string, opts ...grpc.DialOption) (*Client, error) {
	base, err := baseClient.NewClientWithOptions(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor client: %w", err)
	}

	return &Client{
		client: executorpb.NewExecutorServiceClient(base.Connection()),
		base:   base,
	}, nil
}

func (c *Client) Close() error {
	return c.base.Close()
}

func (c *Client) Execute(ctx context.Context, language, code string, problemID int) (*executorpb.ExecuteResponse, error) {
	req := &executorpb.ExecuteRequest{
		Language:  language,
		Code:      code,
		ProblemId: int32(problemID),
	}

	return c.client.Execute(ctx, req)
}

func (c *Client) GetJobStatus(ctx context.Context, jobID string) (*executorpb.GetJobStatusResponse, error) {
	req := &executorpb.GetJobStatusRequest{
		JobId: jobID,
	}

	return c.client.GetJobStatus(ctx, req)
}
