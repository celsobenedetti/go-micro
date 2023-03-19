package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/celso-patiri/go-micro/logger/data"
	"github.com/celso-patiri/go-micro/logger/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (s *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := s.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "Failed"}
		return res, err
	}

	res := &logs.LogResponse{
		Result: "Logged via gRPC!",
	}
	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()
	// register service with autogen code
	logs.RegisterLogServiceServer(s, &LogServer{
		Models: app.Models,
	})

	log.Printf("gRPC server started on port %s", grpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
