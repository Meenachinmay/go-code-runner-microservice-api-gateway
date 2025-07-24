package coding_tests

import (
	"context"
	"fmt"
	codingtestspb "go-code-runner-microservice/api-gateway/go-code-runner-microservice/proto/coding_tests/v1"
	baseClient "go-code-runner-microservice/api-gateway/internal/service/grpc"
)

type Client struct {
	client codingtestspb.CodingTestServiceClient
	base   *baseClient.Client
}

func NewClient(address string) (*Client, error) {
	base, err := baseClient.NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create coding tests client: %w", err)
	}

	return &Client{
		client: codingtestspb.NewCodingTestServiceClient(base.Connection()),
		base:   base,
	}, nil
}

func (c *Client) Close() error {
	return c.base.Close()
}

func (c *Client) VerifyTest(ctx context.Context, testID string) (*codingtestspb.VerifyTestResponse, error) {
	req := &codingtestspb.VerifyTestRequest{
		TestId: testID,
	}

	return c.client.VerifyTest(ctx, req)
}

func (c *Client) StartTest(ctx context.Context, testID, candidateName, candidateEmail string) (*codingtestspb.StartTestResponse, error) {
	req := &codingtestspb.StartTestRequest{
		TestId:         testID,
		CandidateName:  candidateName,
		CandidateEmail: candidateEmail,
	}

	return c.client.StartTest(ctx, req)
}

func (c *Client) SubmitTest(ctx context.Context, testID, code string, passedPercentage int32) (*codingtestspb.SubmitTestResponse, error) {
	req := &codingtestspb.SubmitTestRequest{
		TestId:           testID,
		Code:             code,
		PassedPercentage: passedPercentage,
	}

	return c.client.SubmitTest(ctx, req)
}

func (c *Client) GenerateTest(ctx context.Context, companyID, problemID, expiresInHours int32, clientId string) (*codingtestspb.GenerateTestResponse, error) {
	req := &codingtestspb.GenerateTestRequest{
		CompanyId:      companyID,
		ProblemId:      problemID,
		ExpiresInHours: expiresInHours,
		ClientId:       clientId,
	}

	return c.client.GenerateTest(ctx, req)
}

func (c *Client) GetCompanyTests(ctx context.Context, companyID int32) (*codingtestspb.GetCompanyTestsResponse, error) {
	req := &codingtestspb.GetCompanyTestsRequest{
		CompanyId: companyID,
	}

	return c.client.GetCompanyTests(ctx, req)
}
