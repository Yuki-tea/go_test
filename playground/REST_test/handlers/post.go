package handlers

import (
	"encoding/json"
	"net/http"

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
	defer rows.Close() // CRITICAL: Always close rows when done!

	var posts []BlogPost

	// loop through the results
	for rows.Next() {
		var p BlogPost
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

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	var post BlogPost

	row := db.DB.QueryRow("SELECT id, title, content From blog_posts WHERE id = 1")
	// need to pass the pointers
	// the content of the row will be passed to post
	err := row.Scan(&post.ID, &post.Title, &post.Content)

	if err != nil {
		http.Error(w, "Failed to fetch post from database", http.StatusInternalServerError)
		return
	}
	// tell the browser we are sending JSON, not a plain text
	w.Header().Set("Content-Type", "application/json")
	// encode thte Go struct into JSON and send it
	json.NewEncoder(w).Encode(post)
}
