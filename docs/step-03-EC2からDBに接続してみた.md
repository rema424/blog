# EC2 から DB に接続してみた

## 目次

1. ユーザー・テーブル・初期データを作成する
1. Go プログラムを書く
1. Go プログラムをデプロイする
1. ブラウザからアクセスする

## ユーザー・テーブル・初期データを作成する

SeauelPro で RDS インスタンスに接続してデータベースを選択。

```sql
CREATE DATABASE blog;
SHOW DATABASES;
CREATE USER developer@'%' IDENTIFIED BY 'Passw0rd!';
SELECT host, user FROM mysql.user;
GRANT ALL PRIVILEGES ON blog.* TO developer@'%';
SHOW GRANTS FOR developer@'%';
USE blog;
CREATE TABLE people (
    id int AUTO_INCREMENT,
    name varchar(255),
    age int,
    PRIMARY KEY(id)
);
SHOW TABLES;
DESC people;
INSERT INTO people (name, age) VALUES
('Alice', 21),
('Bob', 22),
('Carol', 23),
('Dave', 24),
('Eve', 25);
```

## Go プログラムを書く

```go
package main

import (
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
	dbx, err := sqlx.Connect("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	dbx.SetMaxOpenConns(30)
	dbx.SetMaxIdleConns(30)
	dbx.SetConnMaxLifetime(60 * time.Second)
	return dbx, nil
}
```

## Go プログラムをデプロイする

```sh
go mod init app
go mod tidy
scp -i ~/Downloads/WebServer1.pem ./playground/main.go ec2-user@$WEB_SERVER_1_IP:~/
scp -i ~/Downloads/WebServer1.pem ./playground/go.mod ec2-user@$WEB_SERVER_1_IP:~/
ssh -i ~/Downloads/WebServer1.pem ec2-user@$$WEB_SERVER_1_IP
```

```sh
export PATH=$PATH:/usr/local/go/bin
go mod tidy
go run main.go
dial tcp 127.0.0.1:3306: connect: connection refused
export DB_USER=developer
export DB_PASS=Passw0rd!
export DB_NAME=blog
export DB_HOST=blog-database.cluster-cll9xfuraffh.ap-northeast-1.rds.amazonaws.com
export DB_PORT=3306
export DB_NET=tcp
export DB_ADDR=''
go run main.go
listen tcp :80: bind: permission denied
go build -o app
sudo ./app # 環境変数が受け継がれない
sudo -E ./app # 環境変数が受け継がれる
```

[sudo 時の環境変数上書き / 引き継ぎについて](https://qiita.com/chroju/items/375582799acd3c5137c7)

## ブラウザからアクセスする

http://IP アドレス/people
http://IP アドレス/people/1
http://IP アドレス/people/5
