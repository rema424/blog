package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

var (
	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")
	DB_NAME = os.Getenv("DB_NAME")
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_NET  = os.Getenv("DB_NET")
	DB_ADDR = os.Getenv("DB_ADDR")
)

func main() {
	fmt.Println(DB_USER)
	fmt.Println(DB_PASS)
	fmt.Println(DB_NAME)
	fmt.Println(DB_HOST)
	fmt.Println(DB_PORT)
	fmt.Println(DB_NET)
	fmt.Println(DB_ADDR)

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(200)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	log.Printf("Listening on localhost:%s", port)
	http.Handle("/", e)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
