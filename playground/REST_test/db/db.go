package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// We import the driver using an underscore. 
	// This tells Go: "Load this package so it installs itself into database/sql, 
	// but I won't call its functions directly."
	_ "github.com/lib/pq"
)

// Capitalize the first letter makes it public scope
var DB *sql.DB

func Init() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	// The Connection String (DSN)
	// the host is "db" - Docker automatically routes this to your Postgres container!
	connStr := fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable", user, password, dbName)
	
	var err error

	// Open the connection
	// If we used ':=', Go would create a brand new LOCAL variable named DB(not the one declared at the top)
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open a DB connection:", err)
	}
	// check the connection
	if err = DB.Ping(); err != nil {
		log.Fatal("DB connection failed: ", err)
	}
	fmt.Println("DB connection successful!")

	initializeTables()
}

func initializeTables() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS blog_posts(
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL
	);`

  // the blank identifier '_' is to ignore the result
	if _, err := DB.Exec(createTableSQL); err != nil {
		log.Fatal("Failed to create table: ", err)
	}

	// insert a test post
	insertPostSQL := `
	INSERT INTO blog_posts (id, title, content)
	VALUES (1, 'My First Go API', 'This data was pulled directly from PostgreSQL!')
	ON CONFLICT (id) DO NOTHING;`
	 
	if _, err := DB.Exec(insertPostSQL); err != nil {
		log.Fatal("Failed to insert dummy data:", err)
	}
	fmt.Println("Database initialized successfully!")
	
}
