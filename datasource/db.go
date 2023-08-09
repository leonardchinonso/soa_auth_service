package datasource

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/leonardchinonso/auth_service_cmp7174/config"
)

// DatabaseContext contains the mongo database pointer
// and the context information for the application
type DatabaseContext struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

// InitDB initializes and initiates connection to the db
// it pings the database to be sure it is up
func InitDB(configMap *map[string]string) (*DatabaseContext, error) {
	uri := fmt.Sprintf("%s%s", (*configMap)[config.BaseUri], (*configMap)[config.DatabaseName])
	client, ctx, cancel, err := connectDatabase(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// ping the database to make sure all connections were successful
	err = ping(client, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	db := client.Database(config.Map[config.DatabaseName])

	return &DatabaseContext{
		Client:     client,
		Database:   db,
		Ctx:        ctx,
		CancelFunc: cancel,
	}, nil
}

// connectDatabase takes in a URI string, connects the mongo server
// to it and returns the mongo client
func connectDatabase(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	// set the deadline for connection initiation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// fetch the client from the mongo connection
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		// if there is an error, cancel the context before returning to avoid context/memory leak
		defer cancel()

		return nil, nil, nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return client, ctx, cancel, nil
}

// ping sends a request through the mongoDB client to ensure a connection has been made
func ping(client *mongo.Client, ctx context.Context) error {
	// mongo.Client has a Ping method to ping mongoDB, the deadline of the Ping method
	// will be determined by ctx. Ping method return error if any occurred, then the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	log.Printf("Connected to database client successfully\n")

	return nil
}
