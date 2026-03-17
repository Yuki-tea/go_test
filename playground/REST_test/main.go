package main

import (
	"os"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	// We import the driver using an underscore. 
	// This tells Go: "Load this package so it installs itself into database/sql, 
	// but I won't call its functions directly."
	_ "github.com/lib/pq"
)
// the fields must start with a capital letter to make it public
// the struct tags show how to translate the data into JSON format
type BlogPost struct {
	ID	int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
}

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

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS blog_posts(
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL
	);`
  // the blank identifier '_' is to ignore the result
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	// insert a test post
	insertPostSQL := `
	INSERT INTO blog_posts (id, title, content)
	VALUES (1, 'My First Go API', 'This data was pulled directly from PostgreSQL!')
	ON CONFLICT (id) DO NOTHING;`
	_, err = db.Exec(insertPostSQL) 
	if err != nil {
		log.Fatal("Failed to insert dummy data:", err)
	}
	fmt.Println("Database initialized successfully!")

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "The database connection was successful!")
	})

	http.HandleFunc("/api/post", func(w http.ResponseWriter, r *http.Request) {
		var post BlogPost
		row := db.QueryRow("SELECT id, title, content From blog_posts WHERE id = 1")
		// need to pass the pointers
		err := row.Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {
			http.Error(w, "Failed to fetch post from database", http.StatusInternalServerError)
			return
		}
		// tell the browser we are sending JSON, not a plain text
		w.Header().Set("Content-Type", "application/json")
		// encode thte Go struct into JSON and send it
		json.NewEncoder(w).Encode(post)
	})

	fmt.Println("Web server is starting on port 8080...")
	// ListenAndServe blocks the program from exiting. We wrap it in log.Fatal 
	// so if the server crashes, it prints the error and exits gracefully.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
