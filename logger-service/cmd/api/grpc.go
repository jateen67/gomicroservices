package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jateen67/log-service/data"
	"github.com/jateen67/log-service/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	// log to mongo
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// return response
	res := &logs.LogResponse{Result: "logged via grpc"}
	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}

	// new grpc server
	s := grpc.NewServer()

	// register the service
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("grpc server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}
}
