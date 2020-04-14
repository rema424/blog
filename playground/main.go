package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

var (
	db *sqlx.DB
)

type Person struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func main() {
	fmt.Println(DB_USER)
	fmt.Println(DB_PASS)
	fmt.Println(DB_NAME)
	fmt.Println(DB_HOST)
	fmt.Println(DB_PORT)
	fmt.Println(DB_NET)
	fmt.Println(DB_ADDR)

	var err error
	db, err = DB()
	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(200)
	})
	e.GET("/people", func(c echo.Context) error {
		var ps []*Person
		if err := db.Select(&ps, "select id, name, age from people;"); err != nil {
			return c.String(500, err.Error())
		}
		return c.JSON(200, ps)
	})
	e.GET("/people/:id", func(c echo.Context) error {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		var p Person
		if err := db.Get(&p, "select id, name, age from people where id = ?;", id); err != nil {
			return c.String(404, err.Error())
		}
		return c.JSON(200, p)
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

func DB() (*sqlx.DB, error) {
	cfg := mysql.Config{
		User:                 DB_USER,
		Passwd:               DB_PASS,
		Net:                  DB_NET,
		Addr:                 DB_ADDR,
		DBName:               DB_NAME,
		Collation:            "utf8mb4_bin",
		InterpolateParams:    true,
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	if cfg.Addr == "" {
		cfg.Addr = DB_HOST + ":" + DB_PORT
	}
	log.Println("[" + DB_ADDR + "]")
	log.Println(cfg.Addr)
	dbx, err := sqlx.Connect("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	dbx.SetMaxOpenConns(30)
	dbx.SetMaxIdleConns(30)
	dbx.SetConnMaxLifetime(60 * time.Second)
	return dbx, nil
}
