package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-pi-proj/getMessage"
	"go-pi-proj/saveMessage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	MongoDBURI = "mongodb://localhost:27017" // Adjust if you have custom settings
)

func main() {

	client, err := InitMongoDB()

	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Connection"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))

	router.POST("/save", func(c *gin.Context) {
		saveMessage.SaveMessage(c, client)
	})

	router.GET("/getMessage", func(c *gin.Context) {
		getMessage.GetMessage(c, client)
	})
	router.Run(":8081") // start the server on port 8081
}

func InitMongoDB() (*mongo.Client, error) {
	// Setup a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use mongo.Connect to establish a connection
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDBURI))
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure the connection is established
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB!")

	return client, nil
}
