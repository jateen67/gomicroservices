package main

import (
	"context"
	"log"
	"time"

	"github.com/jateen67/log-service/data"
)

// any time we want to setup rpc, we need to specify a specific type for it
type RPCServer struct{}

// we also want to define the kind of payload we want to receive from rpc
type RPCPayload struct {
	Name string
	Data string
}

// now we define methods we want to expose via rpc
func (r *RPCServer) LogInfo(payload RPCPayload, res *string) error {
	// write to the logger service, meaning write to mongo
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo:", err)
		return err
	}

	// send our message back to the people who call the method
	*res = "Processed payload via RPC: " + payload.Name + "!"

	return nil
}
