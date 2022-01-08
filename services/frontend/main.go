package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	port = "3000"
)

type frontendServer struct {
	orderSvcAddr string
}

func main() {
	app := gin.Default()

	setupRouter(app)
	fmt.Printf("Starting Frontend service at: %s", port)
	app.Run(":" + port)
}

func setupRouter(app *gin.Engine) {
	svc := new(frontendServer)
	app.GET("/", svc.homeHandler)
	app.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
}
