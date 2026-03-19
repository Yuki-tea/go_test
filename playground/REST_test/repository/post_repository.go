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
