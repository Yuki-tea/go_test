package handlers

import (
	"fmt"
	"encoding/json"
	"net/http"
	"strconv"
	"database/sql"
	"log"

	// custom package
	"rest-api/db"
	"rest-api/models"
	"rest-api/repository"
)

func GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	repo := &repository.PostgresPostRepository{}
	posts, err := repo.GetAll()
	if err != nil {
		log.Printf("GetAllPosts error: %v\n", err)
		http.Error(w, "Failed to fetch posts from the database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	// grab the dynamic {id} from the URL
	idStr := r.PathValue("id")

	// string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	repo := &repository.PostgresPostRepository{}
	post, err := repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound) // 404
			return
		}
		
		// if there was unspecified errors (TCP connection dropped, bad syntax, etc...)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	// tell the browser we are sending JSON, not a plain text
	w.Header().Set("Content-Type", "application/json")
	// encode the Go struct into JSON and send it
	json.NewEncoder(w).Encode(post)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var newPost models.BlogPost
	
	// read the incoming JSON body and decode it into our struct
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	repo := &repository.PostgresPostRepository{}
	err = repo.Create(&newPost)
	if err != nil {
		http.Error(w, "Failed to save to database", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	// 201 Created
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// grab the dynamic {id} from the URL
	idStr := r.PathValue("id")

	// string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	// Exec for SQL commands with no return rows
	result, err := db.DB.Exec("DELETE FROM blog_posts WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	// check if the post actually existed
	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func PutPostHandler(w http.ResponseWriter, r *http.Request) {
	// grab the dynamic {id} from the URL
	idStr := r.PathValue("id")

	// string to integer
	// need this to handle multiple data types
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedPost models.BlogPost
	// read the incoming JSON body and decode it into our struct
	err = json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	
	putSQL := `
		UPDATE blog_posts SET title = $1, content = $2 WHERE id = $3 RETURNING id
	`

	err = db.DB.QueryRow(putSQL, updatedPost.Title, updatedPost.Content, id).Scan(&updatedPost.ID)
	if err != nil {
		http.Error(w, "Failed to update database", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	// 200 OK
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPost)
}

func PatchPostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// empty interface
	// key = string & value = anything
	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	if len(updates) == 0 {
		http.Error(w, "No fields provided to update", http.StatusBadRequest)
		return
	}

	// dynamically build the SQL query
	query := "UPDATE blog_posts SET "
	// empty array
	// anything can be hold
	// need this to handle multiple data types
	args := []interface{}{} // Holds the actual values we will inject
	argId := 1 // keeps track of our $1, $2 placeholders

	if title, exists := updates["title"]; exists {
		query += fmt.Sprintf("title = $%d, ", argId)
		args = append(args, title)
		argId++
	}

	if content, exists := updates["content"]; exists {
		query += fmt.Sprintf("content = $%d, ", argId)
		args = append(args, content)
		argId++
	}

	// slice off the trailing comma and space from our loop
	query = query[:len(query)-2]

	// append the WHERE clause and returning fields
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, title, content", argId)
	args = append(args, id)

	var updatedPost models.BlogPost
	err = db.DB.QueryRow(query, args...).Scan(&updatedPost.ID, &updatedPost.Title, &updatedPost.Content)
	if err != nil {
		http.Error(w, "Post not found or failed to update", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	// 200 OK
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPost)
}
