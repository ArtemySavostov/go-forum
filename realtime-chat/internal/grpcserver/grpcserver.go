package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/ArtemySavostov/chat-protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedChatServiceServer
	rooms map[string]*Room
	mu    sync.Mutex
}

type Room struct {
	id      string
	Clients map[string]chan *pb.ChatMessage
}

func NewServer() *Server {
	return &Server{
		rooms: make(map[string]*Room),
	}
}

func (s *Server) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	roomID := generateRoomID()
	room := &Room{
		id:      roomID,
		Clients: make(map[string]chan *pb.ChatMessage),
	}
	s.rooms[roomID] = room
	log.Printf("Room created: %s by client %s", roomID, req.ClientId)

	return &pb.CreateRoomResponse{RoomId: roomID, Success: true}, nil
}

func (s *Server) SubscribeToRoom(req *pb.SubscribeToRoomRequest, stream pb.ChatService_SubscribeToRoomServer) error {
	roomID := req.RoomId
	clientID := req.ClientId
	s.mu.Lock()
	room, ok := s.rooms[roomID]
	if !ok {
		return fmt.Errorf("room not found: %s", roomID)
	}
	messageChan := make(chan *pb.ChatMessage, 10)
	room.Clients[clientID] = messageChan
	s.mu.Unlock()

	log.Printf("Client %s subscribed to room %s", clientID, roomID)

	defer func() {
		delete(room.Clients, clientID)
		close(messageChan)
		s.mu.Unlock()
		log.Printf("Client %s unsubscribed from room %s", clientID, roomID)
	}()
	for {
		select {
		case message := <-messageChan:
			err := stream.Send(message)
			log.Printf("Error sending message to client %s: %v", clientID, err)
			return err
		case <-stream.Context().Done():
			return nil
		}

	}
}

// func (s *Server) RelayMessage(ctx context.Context, req *pb.RelayMessageRequest) (*pb.Empty, error) {
// 	// Create the full message (you might need more information here)
// 	fullMessage := []byte(req.RoomId + ": " + req.Message)

// 	// Send the full message to the broadcast channel
// 	s.hub.Broadcast <- fullMessage

//		return &pb.Empty{}, nil
//	}
func (s *Server) SendMessage(ctx context.Context, req *pb.ChatMessage) (*pb.SendMessageResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	room, ok := s.rooms[req.RoomId]
	if !ok {
		return &pb.SendMessageResponse{Success: false, Error: "Room not found"}, nil
	}
	message := &pb.ChatMessage{
		RoomId:    req.RoomId,
		ClientId:  req.ClientId,
		Text:      req.Text,
		Timestamp: time.Now().Unix(),
	}
	for _, messageChan := range room.Clients {
		messageChan <- message
	}
	log.Printf("Message sent to room %s from client %s: %s", req.RoomId, req.ClientId, req.Text)
	return &pb.SendMessageResponse{Success: true}, nil
}

func generateRoomID() string {
	return fmt.Sprintf("room-%d", time.Now().UnixNano())
}
func StartGRPCServer(lis net.Listener) error {
	s := grpc.NewServer()
	chatServer := NewServer()
	pb.RegisterChatServiceServer(s, chatServer)

	reflection.Register(s)

	log.Printf("gRPC server listening on %v", lis.Addr())
	return s.Serve(lis)
}
