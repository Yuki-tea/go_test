package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	// custom package
	"rest-api/db"
)

// the fields must start with a capital letter to make it public
// the struct tags show how to translate the data into JSON format
type BlogPost struct {
	ID	int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
}

func GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, content FROM blog_posts ORDER BY id ASC")
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	// if you forget this, the TCP connection remains active
	defer rows.Close() // CRITICAL: Always close rows when done!

	var posts []BlogPost

	// loop through the results (like a Iterator in Java)
	for rows.Next() {
		var p BlogPost
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
		posts = []BlogPost{}
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

	var post BlogPost

	row := db.DB.QueryRow("SELECT id, title, content From blog_posts WHERE id = $1", id)
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
	var newPost BlogPost
	
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
