package repository

import (
	"fmt"
	"errors"

	"database/sql"
	"rest-api/db"
	"rest-api/models"
)

// the contract(interface)
// Any database (Postgres, MongoDB, or a fake one for testing) 
// MUST implement these exact methods to be allowed in our app.
type PostRepository interface {
	GetAll() ([]models.BlogPost, error)
	GetByID(id int) (models.BlogPost,error)
	Create(post *models.BlogPost) error
	Update(id int, post models.BlogPost) (models.BlogPost, error)
	Patch(id int, updates map[string]interface{}) (models.BlogPost, error)
	Delete(id int) error 
}

// represent the specific PostGreSQL DB
type PostgresPostRepository struct {}

// link the function to the struct, like OOP paradigm
// r behaves like "this" in Java and TS (you can decide the name by yourself)
// 1-to-2 letters are preffered
func(r *PostgresPostRepository) GetAll() ([]models.BlogPost, error) {
	rows, err := db.DB.Query("SELECT id, title, content FROM blog_posts ORDER BY id ASC")
	if err != nil {
		return nil, err // return standard errors here, not HTTP errors
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
			return nil, err
		}
		posts = append(posts, p)
	}

	// if none, return an empty array
	if posts == nil {
		posts = []models.BlogPost{}
	}
	return posts, nil
}

// you have to use r here, because it's used as r in GetAll func above
func (r *PostgresPostRepository) GetByID(id int) (models.BlogPost, error) {
	var post models.BlogPost

	row := db.DB.QueryRow("SELECT id, title, content FROM blog_posts WHERE id = $1", id)
	// need to pass the pointers
	// the content of the row will be passed to post
	err := row.Scan(&post.ID, &post.Title, &post.Content)

	if err != nil {
		return post, err
	}

	return post, nil
}

func (r *PostgresPostRepository) Create(post *models.BlogPost) error {
	// use $1 and $2 as placeholders to prevent SQL Injection attacks!
	// returning auto-generated ID
	insertSQL := `
		INSERT INTO blog_posts (title, content) VALUES ($1, $2) RETURNING id
	`
	err := db.DB.QueryRow(insertSQL, post.Title, post.Content).Scan(&post.ID)
	if err != nil {
		return err
	}
	return nil
}
