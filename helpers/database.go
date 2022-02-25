package helpers

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, context.Context) {
	cnnStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority",
		GetSetting("DATABASE.USER"),
		GetSetting("DATABASE.PASSWORD"),
		GetSetting("DATABASE.URI"),
		GetSetting("DATABASE.NAME"))

	client, err := mongo.NewClient(options.Client().ApplyURI(cnnStr))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client, ctx
}
