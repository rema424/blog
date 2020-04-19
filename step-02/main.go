package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	http.Handle("/", e)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(c echo.Context) error {
	if c.Request().URL.Path != "/" {
		return c.String(http.StatusNotFound, "404 page not found")
	}
	return c.String(http.StatusOK, "Hello, World!")
}
