package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"../eventpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := eventpb.NewEventServiceClient(cc)
	//fmt.Printf("Created client: %f", c)
	//doUnary(c)
	//doStreaming(c)
	biStream(c)
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

func doStreaming(c eventpb.EventServiceClient) {
	req := &eventpb.EventManyTimesRequest{
		Eventing: &eventpb.Eventing{
			EventName:   "ALLO DENIZ",
			Description: "SUPER DESCRIPTION",
		},
	}

	resStream, err := c.EventManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Event Stream RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Printf("Response from stream: %v", msg.GetResult())
	}
}

func biStream(c eventpb.EventServiceClient) {
	fmt.Println("Starting to do a Bi Stream RPC")
	requests := []*eventpb.EventBiRequest{
		&eventpb.EventBiRequest{
			Eventing: &eventpb.Eventing{
				EventName:   "ALLO DENIZ1",
				Description: "SUPER DESCRIPTION",
			},
		},
		&eventpb.EventBiRequest{
			Eventing: &eventpb.Eventing{
				EventName:   "ALLO DENIZ2",
				Description: "SUPER DESCRIPTION",
			},
		},
		&eventpb.EventBiRequest{
			Eventing: &eventpb.Eventing{
				EventName:   "ALLO DENIZ3",
				Description: "SUPER DESCRIPTION",
			},
		},
	}

	stream, err := c.EventBi(context.Background())
	if err != nil {
		log.Fatalf("Error while create stream %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() { //send msg
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(100 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received message: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc

}
