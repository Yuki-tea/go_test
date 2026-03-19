package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"rest-api/models"
)

// fake DB
type MockPostRepository struct {
	MockData []models.BlogPost
}
// test functions in Go must always start with the word "Test" and take (t *testing.T)
func TestGetAllPostsHandler(t *testing.T) {
	fakeRepo := &MockPostRepository{
		MockData: []models.BlogPost{
			{ID: 1, Title: "Test Driven Go", Content: "Testing is fun?"},
		},
	}
	
	// inject the fake database into the Handler struct
	handler := &PostHandler{
		Repo: fakeRepo,
	}

	// fake HTTP Request
	req, err := http.NewRequest("GET", "/api/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	// fake ResponseRecorder to capture what the handler sends back
	rr := httptest.NewRecorder()
	// execute the handler directly
	handler.GetAllPostsHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	t.Logf("Response Body: %v", rr.Body.String())
}

func (m *MockPostRepository) GetAll() ([]models.BlogPost, error) {
	return m.MockData, nil
}

func (m *MockPostRepository) GetByID(id int) (models.BlogPost, error) { return models.BlogPost{}, nil }
func (m *MockPostRepository) Create(post *models.BlogPost) error { return nil }
func (m *MockPostRepository) Update(id int, post models.BlogPost) (models.BlogPost, error) { return models.BlogPost{}, nil }
func (m *MockPostRepository) Patch(id int, updates map[string]interface{}) (models.BlogPost, error) { return models.BlogPost{}, nil }
func (m *MockPostRepository) Delete(id int) error { return nil }
