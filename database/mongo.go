package database

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB(url string) *mongo.Client {
	options := options.Client().ApplyURI(url)
	options.SetHTTPClient(newHTTPClient())
	options.SetConnectTimeout(mongoConnectTimeout)
	options.SetMaxConnIdleTime(connMaxIdleTime)
	options.SetMaxConnecting(0) // use default = 2
	options.SetSocketTimeout(mongoSocketTimeout)
	options.SetTimeout(mongoTransactionTimeout)
	options.SetMaxPoolSize(mongoMaxPoolSize)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options)
	if err != nil {
		log.Panic("error while creating connection to the database!!", err)
	}

	// Check the connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Panic("could not ping database", err)
	}

	return client
}

func IsMongoReady() bool {
	// TODO: implement check if mongo is ready
	return true
}

func newHTTPClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = mongoMaxPoolSize
	t.MaxConnsPerHost = mongoMaxPoolSize
	t.MaxIdleConnsPerHost = mongoMaxPoolSize
	return &http.Client{
		Timeout:   mongoTransactionTimeout,
		Transport: t,
	}
}
