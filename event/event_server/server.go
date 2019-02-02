package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"../eventpb"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct{}

type EventItem struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	EventId     string             `bson:"id"`
	EventName   string             `bson:"event_name"`
	Description string             `bson:"description"`
}

func (*server) Event(ctx context.Context, req *eventpb.EventRequest) (*eventpb.EventResponse, error) {
	eventName := req.GetEventing().GetEventName()
	eventDescription := req.GetEventing().GetDescription()
	result := "Event name " + eventName + " Description " + eventDescription
	res := &eventpb.EventResponse{
		Result: result,
	}

	return res, nil
}

func (*server) EventManyTimes(req *eventpb.EventManyTimesRequest, stream eventpb.EventService_EventManyTimesServer) error {
	eventName := req.GetEventing().GetEventName()
	eventDescription := req.GetEventing().GetDescription()

	result := "Event name " + eventName + " Description " + eventDescription
	res := &eventpb.EventManyTimesResponse{
		Result: result,
	}
	stream.Send(res)

	return nil
}

func (*server) EventBi(stream eventpb.EventService_EventBiServer) error {
	fmt.Printf("Bi Event function invoke")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error Deniz Fault %v", err)
			return err
		}


		data := EventItem{
			EventId:     req.GetEventing().GetEventId(),
			EventName:   req.GetEventing().GetEventName(),
			Description: req.GetEventing().GetDescription(),
		}

		res, err := collection.InsertOne(context.Background(), data)
		if err != nil {

			log.Fatalf("Internal error: %v", err)
			
		}

		oid, ok := res.InsertedID.(primitive.ObjectID)
		if !ok {

				fmt.Println("Cannot convert to OID")
			
		}
			fmt.Println("Id " + oid.Hex())
			fmt.Println("EventId " +req.GetEventing().GetEventId())
			fmt.Println("EventName " +req.GetEventing().GetEventName())
			fmt.Println("Description " +req.GetEventing().GetDescription())
		

       fmt.Println("No Error Inserrt")
	}
	
}

func main() {
	fmt.Println("Hello World")

	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database("BFitOps").Collection("Event")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	eventpb.RegisterEventServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	client.Disconnect(context.TODO())
}
