package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	port = "8080"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello Andrew!!")
	})

	r.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})

	return r
}

func main() {
	fmt.Printf("Starting Order service at: %s", port)
	r := setupRouter()
	r.Run(":" + port)
}
