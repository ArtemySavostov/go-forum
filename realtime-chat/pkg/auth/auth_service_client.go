package auth

import (
	"context"
	"fmt"
	"log"

	pb "realtime-chat/pkg/auth/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient interface {
	ValidateToken(ctx context.Context, token string) (string, error)
}

type authServiceClient struct {
	client pb.AuthServiceClient
}

func NewAuthServiceClient(addr string) (AuthServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	c := pb.NewAuthServiceClient(conn)

	return &authServiceClient{client: c}, nil
}

func (c *authServiceClient) ValidateToken(ctx context.Context, token string) (string, error) {
	resp, err := c.client.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: token})
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %w", err)
	}
	return resp.UserId, nil
}
