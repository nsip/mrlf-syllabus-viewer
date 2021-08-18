package main

import (
	// "net/http"

	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/v4/middleware"

	"log"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!\n")
	// })

	e.Static("/", "public")

	log.Println("==================================================")
	log.Println("Static content server running...")
	log.Println("browse to http://localhost:1323/")
	log.Println("content will be served from the /public directory")
	log.Println("Ctrl-C to shut down server.")
	log.Println("==================================================")

	// Start server
	// e.Run(standard.New(":8080"))
	e.Logger.Fatal(e.Start(":1323"))

}
