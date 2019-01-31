package main

import (
	"context"
	"fmt"
	"log"

	"../eventpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: v%", err)
	}

	defer cc.Close()

	c := eventpb.NewEventServiceClient(cc)
	//fmt.Printf("Created client: %f", c)
	doUnary(c)
}

func doUnary(c eventpb.EventServiceClient) {
	req := &eventpb.EventRequest{
		Eventing: &eventpb.Eventing{
			EventName:   "ALLO DENIZ",
			Description: "SUPER DESCRIPTION",
		},
	}
	res, err := c.Event(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Event RPC: %v", err)
	}

	log.Printf("Response from Event: %v", res.Result)
}
