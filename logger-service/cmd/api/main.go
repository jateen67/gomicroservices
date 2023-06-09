package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/jateen67/log-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	port     = "80"
	rpcPort  = "5001"
	gRpcPort = "50001"
	mongoUrl = "mongodb://mongo:27017" // 'mongo' is the name defined for our mongo instance in our docker-compose file
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// register the rpc server to tell the app that we will indeed be accepting rpc requests
	err = rpc.Register(new(RPCServer))
	// start our rpc server
	go app.rpcListen()

	// start our grpc server
	go app.gRPCListen()

	// start our server
	log.Printf("starting logger service on port %s\n", port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

// method to start up rpc server
func (app *Config) rpcListen() error {
	log.Println("starting rpc server on port ", rpcPort)
	// listen for a tcp connection on any and all ip addresses (0.0.0.0) and our rpcPort
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	// loop that executes forever to accept connections
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoUrl)

	// set the authentication
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error connecting:", err)
		return nil, err
	}

	log.Println("connected to mongo")

	return c, nil
}
