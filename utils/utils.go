package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DatabaseConnect(dbname string) (db *mongo.Database, err error) {
	/*
		commandMonitor := &event.CommandMonitor{
			Started: func(ctx context.Context, event *event.CommandStartedEvent) {
				fmt.Printf("Executing command: %s\nCommand details: %s\n", event.CommandName, event.Command.String())
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(commandMonitor)
		client, err := mongo.Connect(ctx, clientOptions.ApplyURI("mongodb://localhost:27017"))
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println("database connect error", err)
		return nil, err
	}
	db = client.Database(dbname)
	return db, nil
}
