package saveMessage

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
)

type Message struct {
	Message []string `json:"message"`
}

func SaveMessage(c *gin.Context) {
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

	db, err := sql.Open("sqlite3", "/path/to/your/database/file.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to the database",
		})
		return
	}
	defer db.Close()

	for _, char := range message.Message {
		_, err := db.Exec("INSERT INTO message (character) VALUES (?)", char)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to insert into the database",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
	})
}
