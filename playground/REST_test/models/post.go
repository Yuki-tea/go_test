package models

// this is the absolute core of the app
// no HTTP or PostgreSQL even exist here
type BlogPost struct {
	// the struct tags show how to translate the data into JSON format
	ID	int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
}
