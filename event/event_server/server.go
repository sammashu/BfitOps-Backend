package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"../eventpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Event(ctx context.Context, req *eventpb.EventRequest) (*eventpb.EventResponse, error) {
	eventName := req.GetEventing().GetEventName()
	eventDescription := req.GetEventing().GetDescription()
	result := "Event name " + eventName + " Description " + eventDescription
	res := &eventpb.EventResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	eventpb.RegisterEventServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
