package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "8081"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

var counts int64

func main() {
	log.Println("Statrting authentication service")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Cant connect to Postgres")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := "host=localhost user=postgres password=super369 dbname=users port=5432 sslmode=disable TimeZone=Asia/Kolkata connect_timeout=5"
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready...")
			counts++
		} else {
			log.Println("Connected to database!")
			return connection
		}

		if counts > 20 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
