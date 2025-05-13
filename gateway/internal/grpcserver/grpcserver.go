package grpcclient

import (
	"context"
	"fmt"
	"log"

	pb "github.com/ArtemySavostov/chat-protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatServiceClientWrapper struct {
	client pb.ChatServiceClient
	conn   *grpc.ClientConn
}

func NewChatServiceClientWrapper(chatServiceEndpoint string) (*ChatServiceClientWrapper, error) {

	conn, err := grpc.Dial(chatServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to chat service: %v", err)
		return nil, fmt.Errorf("failed to connect to chat service: %w", err)
	}

	client := pb.NewChatServiceClient(conn)

	return &ChatServiceClientWrapper{client: client, conn: conn}, nil
}

func (c *ChatServiceClientWrapper) Close() error {
	return c.conn.Close()
}

func (c *ChatServiceClientWrapper) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	return c.client.CreateRoom(ctx, req)
}

func (c *ChatServiceClientWrapper) SendMessage(ctx context.Context, req *pb.ChatMessage) (*pb.SendMessageResponse, error) {
	return c.client.SendMessage(ctx, req)
}

func (c *ChatServiceClientWrapper) SubscribeToRoom(ctx context.Context, req *pb.SubscribeToRoomRequest) (pb.ChatService_SubscribeToRoomClient, error) {
	return c.client.SubscribeToRoom(ctx)
}

// func HandleCreateRoom(w http.ResponseWriter, r *http.Request, client *ChatServiceClientWrapper) {
// 	fmt.Println("Calling CreateRoom gRPC Service")

// }

// func HandleSendMessage(w http.ResponseWriter, r *http.Request, client *ChatServiceClientWrapper) {
// 	fmt.Println("Calling SendMessage gRPC Service")

// }

// func HandleSubscribeToRoom(w http.ResponseWriter, r *http.Request, client *ChatServiceClientWrapper) {
// 	fmt.Println("Calling SubscribeToRoom gRPC Service")
// }
