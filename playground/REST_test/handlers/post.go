package handlers

import (
	"fmt"
	"encoding/json"
	"net/http"
	"strconv"

	// custom package
	"rest-api/db"
	"rest-api/models"
)

func GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, content FROM blog_posts ORDER BY id ASC")
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	// if you forget this, the TCP connection remains active
	defer rows.Close() // CRITICAL: Always close rows when done!

	var posts []models.BlogPost

	// loop through the results (like a Iterator in Java)
	for rows.Next() {
		var p models.BlogPost
		// rows actually pointing to a single row at each moment
		// and passing the data to the p here
		if err := rows.Scan(&p.ID, &p.Title, &p.Content); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	// if none, return an empty array
	if posts == nil {
		posts = []models.BlogPost{}
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

	var post models.BlogPost

	row := db.DB.QueryRow("SELECT id, title, content FROM blog_posts WHERE id = $1", id)
	// need to pass the pointers
	// the content of the row will be passed to post
	err = row.Scan(&post.ID, &post.Title, &post.Content)

	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
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

	// use $1 and $2 as placeholders to prevent SQL Injection attacks!
	// returning auto-generated ID
	insertSQL := `
		INSERT INTO blog_posts (title, content) VALUES ($1, $2) RETURNING id
	`
	err = db.DB.QueryRow(insertSQL, newPost.Title, newPost.Content).Scan(&newPost.ID)
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
