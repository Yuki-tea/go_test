package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	// custom packages
	"rest-api/db"
)
// the fields must start with a capital letter to make it public
// the struct tags show how to translate the data into JSON format
type BlogPost struct {
	ID	int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
}

func main() {
	db.Init()
	defer db.DB.Close()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "The database connection was successful!")
	})

	http.HandleFunc("/api/post", func(w http.ResponseWriter, r *http.Request) {
		var post BlogPost
		row := db.DB.QueryRow("SELECT id, title, content From blog_posts WHERE id = 1")
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
