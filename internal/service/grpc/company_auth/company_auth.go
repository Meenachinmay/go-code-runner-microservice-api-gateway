package company_auth

import (
	"context"
	"fmt"
	companyauthpb "go-code-runner-microservice/api-gateway/go-code-runner-microservice/proto/company_auth/v1"
	baseClient "go-code-runner-microservice/api-gateway/internal/service/grpc"
)

type Client struct {
	client companyauthpb.CompanyAuthServiceClient
	base   *baseClient.Client
}

func NewClient(address string) (*Client, error) {
	base, err := baseClient.NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create company auth client: %w", err)
	}

	return &Client{
		client: companyauthpb.NewCompanyAuthServiceClient(base.Connection()),
		base:   base,
	}, nil
}

func (c *Client) Close() error {
	return c.base.Close()
}

func (c *Client) Register(ctx context.Context, name, email, password string) (*companyauthpb.RegisterResponse, error) {
	req := &companyauthpb.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	}

	return c.client.Register(ctx, req)
}

func (c *Client) Login(ctx context.Context, email, password string) (*companyauthpb.LoginResponse, error) {
	req := &companyauthpb.LoginRequest{
		Email:    email,
		Password: password,
	}

	return c.client.Login(ctx, req)
}

func (c *Client) GenerateAPIKey(ctx context.Context, companyID int32) (*companyauthpb.GenerateAPIKeyResponse, error) {
	req := &companyauthpb.GenerateAPIKeyRequest{
		CompanyId: companyID,
	}

	return c.client.GenerateAPIKey(ctx, req)
}

func (c *Client) GenerateClientID(ctx context.Context, companyID int32) (*companyauthpb.GenerateClientIDResponse, error) {
	req := &companyauthpb.GenerateClientIDRequest{
		CompanyId: companyID,
	}

	return c.client.GenerateClientID(ctx, req)
}