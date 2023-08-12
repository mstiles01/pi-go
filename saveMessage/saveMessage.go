package saveMessage

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"time"
)

type Message struct {
	Message []string `json:"message"`
}

func SaveMessage(c *gin.Context, client *mongo.Client) {
	var message Message

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error reading request body",
		})
		return
	}

	err = json.Unmarshal(body, &message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Save to MongoDB
	collection := client.Database("splitflap").Collection("message")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, bson.M{"message": message.Message})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save message to database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
	})
}
