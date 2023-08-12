package getMessage

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Message []string `json:"message" bson:"message"`
}

func GetMessage(c *gin.Context, client *mongo.Client) {
	collection := client.Database("splitflap").Collection("message")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var message Message

	// Here we're using FindOne to get the most recent message.
	// We're using options to sort by ID descending to get the most recent entry.
	err := collection.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "No documents found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the message"})
		return
	}

	c.JSON(http.StatusOK, message)
}
