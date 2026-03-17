package main

import (
	"os"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	// We import the driver using an underscore. 
	// This tells Go: "Load this package so it installs itself into database/sql, 
	// but I won't call its functions directly."
	_ "github.com/lib/pq"
)

func main() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	// The Connection String (DSN)
	// the host is "db" - Docker automatically routes this to your Postgres container!
	connStr := fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable", user, password, dbName)

	// Open the connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open a DB connection:", err)
	}
	defer db.Close() // ensures the DB closes right before the main function ends!
	// check the connection
	err = db.Ping()
	if err != nil {
		fmt.Println("connection failed!")
	} else {
		fmt.Println("connection success!")
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "The database connection was successful!")
	})
	fmt.Println("Web server is starting on port 8080...")
	// ListenAndServe blocks the program from exiting. We wrap it in log.Fatal 
	// so if the server crashes, it prints the error and exits gracefully.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
