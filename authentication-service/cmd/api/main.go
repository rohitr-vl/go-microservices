package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

type Config struct {
	Repo data.Repository
}

var counts int64

func main() {
	log.Println("Starting authentication service on port:", webPort)

	conn := connectToDB()
	if conn == nil {
		log.Panic("Cant connect to Postgres")
	}

	app := Config{

	}
	defer conn.Close()
/*
	var dbName, dbVersion, dbUser, dbCatalog, dbTable string
	conn.QueryRow("SELECT current_database()").Scan(&dbName)
	conn.QueryRow("SELECT version()").Scan(&dbVersion)
	conn.QueryRow("SELECT current_user").Scan(&dbUser)
	conn.QueryRow("SELECT * FROM INFORMATION_SCHEMA.INFORMATION_SCHEMA_CATALOG_NAME").Scan(&dbCatalog)
	conn.QueryRow("SELECT * FROM information_schema.tables WHERE table_name LIKE 'users'").Scan(&dbTable)
	log.Printf("\n Current DB: %s, version: %s, user: %s, catalog: %s, table: %s \n",dbName, dbVersion, dbUser, dbCatalog, dbTable)
*/
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
	// connect to postgres
	dsn := os.Getenv("DSN")
	fmt.Println("DSN: ",dsn)
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready...")
			counts++
		} else {
			log.Println("Connected to database!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}