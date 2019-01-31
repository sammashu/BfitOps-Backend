protoc event/eventpb/event.proto --go_out=plugins=grpc:.

go run event/event_server/server.go