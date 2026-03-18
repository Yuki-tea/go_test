package main

import (
	"fmt"
	"log"
	"net/http"

	// custom packages
	"rest-api/db"
	"rest-api/handlers"
)

func main() {
	db.Init()
	defer db.DB.Close()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "The database connection was successful!")
	})

	http.HandleFunc("GET /api/posts", handlers.GetAllPostsHandler)
	http.HandleFunc("GET /api/posts/{id}", handlers.GetPostByIDHandler)
	http.HandleFunc("POST /api/posts", handlers.CreatePostHandler)

	fmt.Println("Web server is starting on port 8080...")
	// ListenAndServe blocks the program from exiting. We wrap it in log.Fatal 
	// so if the server crashes, it prints the error and exits gracefully.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
