package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/rs/cors"

	// custom packages
	"rest-api/db"
	"rest-api/handlers"
	"rest-api/repository"
)

func main() {
	db.Init()
	defer db.DB.Close()

	postgresRepo := &repository.PostgresPostRepository{}
	// inject the repository into the Handler struct
	postHandler := &handlers.PostHandler{
		Repo: postgresRepo,
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "The database connection was successful!")
	})

	// the methods are attached to the postHandler struct
	http.HandleFunc("GET /api/posts", postHandler.GetAllPostsHandler)
	http.HandleFunc("GET /api/posts/{id}", postHandler.GetPostByIDHandler)
	http.HandleFunc("POST /api/posts", postHandler.CreatePostHandler)
	http.HandleFunc("DELETE /api/posts/{id}", postHandler.DeletePostHandler)
	http.HandleFunc("PUT /api/posts/{id}", postHandler.PutPostHandler)
	http.HandleFunc("PATCH /api/posts/{id}", postHandler.PatchPostHandler)

	allowedOrigin := os.Getenv("FRONTEND_URL")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}
	c := cors.New(cors.Options{
		// explicitly allow the Next.js default port
		AllowedOrigins: []string{allowedOrigin},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})
	
	// wrap default router in the CORS 
	protectedHandler := c.Handler(http.DefaultServeMux)

	fmt.Println("Web server is starting on port 8080...")

	// ListenAndServe blocks the program from exiting. We wrap it in log.Fatal 
	// so if the server crashes, it prints the error and exits gracefully.
	log.Fatal(http.ListenAndServe(":8080", protectedHandler))
}
