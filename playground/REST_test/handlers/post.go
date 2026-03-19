package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"database/sql"
	"log"

	// custom package
	"rest-api/models"
	"rest-api/repository"
)

type PostHandler struct {
	Repo repository.PostRepository // not the Postgres struct
}

// add the receiver
func (h *PostHandler) GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Repo.GetAll()
	if err != nil {
		log.Printf("GetAllPosts error: %v\n", err)
		http.Error(w, "Failed to fetch posts from the database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	// grab the dynamic {id} from the URL
	idStr := r.PathValue("id")

	// string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	post, err := h.Repo.GetByID(id)
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

func (h *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var newPost models.BlogPost
	
	// read the incoming JSON body and decode it into our struct
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	err = h.Repo.Create(&newPost)
	if err != nil {
		http.Error(w, "Failed to save to database", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	// 201 Created
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

func (h *PostHandler) PutPostHandler(w http.ResponseWriter, r *http.Request) {
	// grab the dynamic {id} from the URL
	idStr := r.PathValue("id")

	// string to integer
	// need this to handle multiple data types
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var newPost models.BlogPost
	// read the incoming JSON body and decode it into our struct
	err = json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	updatedPost, err := h.Repo.Update(id, newPost)
	if err != nil {
		http.Error(w, "Failed to update database", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	// 200 OK
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPost)
}

func (h *PostHandler) PatchPostHandler(w http.ResponseWriter, r *http.Request) {
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
	updatedPost, err := h.Repo.Patch(id, updates)

	if err != nil {
		http.Error(w, "Post not found or failed to update", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// 200 OK
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPost)
}

func (h *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// grab the dynamic {id} from the URL
	idStr := r.PathValue("id")

	// string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.Repo.Delete(id)
	if err != nil {
		if err.Error() == "post not found" {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
