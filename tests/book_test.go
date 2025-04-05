package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/foverokavindz/book-api/handlers"
	"github.com/foverokavindz/book-api/models"
	"github.com/foverokavindz/book-api/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// setupTestBookStore creates a test store with some sample books
func setupTestBookStore() *storage.BookStore {
	store := storage.NewBookStore("test-storage")
	
	// Add sample books
	store.AddBook(models.Book{
		BookID:          "1",
		Title:           "Test Book 1",
		AuthorID:        "auth-001",
		PublisherID:     "pub-001",
		PublicationDate: "2023-01-15",
		Description:     "Description for test book 1",
		Price:           19.99,
		ISBN:            "1234567890",
		Pages:           320,
		Genre:           "Fiction",
		Quantity:        10,
	})
	
	store.AddBook(models.Book{
		BookID:          "2",
		Title:           "Test Book 2",
		AuthorID:        "auth-002",
		PublisherID:     "pub-002",
		PublicationDate: "2023-02-20",
		Description:     "Description for test book 2",
		Price:           29.99,
		ISBN:            "0987654321",
		Pages:           250,
		Genre:           "Non-Fiction",
		Quantity:        15,
	})
	
	return store
}

func TestGetBooks(t *testing.T) {
	// Setup
	test_store := setupTestBookStore()
	handler := handlers.NewBookHandler(test_store)
	
	// Create a request
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.GetBooks(rr, req)
	
	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code)
	
	// Parse response
	var books []models.Book
	err = json.Unmarshal(rr.Body.Bytes(), &books)
	assert.NoError(t, err)
	
	// Verify response
	assert.Equal(t, 2, len(books))
	assert.Equal(t, "Test Book 1", books[0].Title)
	assert.Equal(t, "Test Book 2", books[1].Title)
}

func TestGetBook(t *testing.T) {
	// Setup
	store := setupTestBookStore()
	handler := handlers.NewBookHandler(store)
	
	// Test cases
	tests := []struct {
		name       string
		bookID     string
		wantStatus int
		wantBook   bool
	}{
		{"Valid ID", "1", http.StatusOK, true},
		{"Invalid ID", "999", http.StatusNotFound, false},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("GET", "/books/"+tc.bookID, nil)
			if err != nil {
				t.Fatal(err)
			}
			
			// Add URL parameters to the request
			vars := map[string]string{
				"id": tc.bookID,
			}
			req = mux.SetURLVars(req, vars)
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handler.GetBook(rr, req)
			
			// Check status code
			assert.Equal(t, tc.wantStatus, rr.Code)
			
			// For successful requests, check the book data
			if tc.wantBook {
				var book models.Book
				err = json.Unmarshal(rr.Body.Bytes(), &book)
				assert.NoError(t, err)
				assert.Equal(t, tc.bookID, book.BookID)
			}
		})
	}
}

func TestCreateBook(t *testing.T) {
	// Setup
	store := setupTestBookStore()
	handler := handlers.NewBookHandler(store)
	
	// Test cases
	tests := []struct {
		name       string
		book       models.Book
		wantStatus int
	}{
		{
			"New Book",
			models.Book{
				Title:           "New Test Book",
				AuthorID:        "auth-003",
				PublisherID:     "pub-003",
				PublicationDate: "2023-03-25",
				Description:     "New Description",
				Price:           39.99,
				ISBN:            "1122334455",
				Pages:           400,
				Genre:           "Mystery",
				Quantity:        5,
			},
			http.StatusCreated,
		},
		{
			"Duplicate Book",
			models.Book{
				BookID:          "1", // Already exists
				Title:           "Duplicate Book",
				AuthorID:        "auth-004",
				PublisherID:     "pub-004",
				Description:     "Description",
				Price:           10.99,
				ISBN:            "9876543210",
			},
			http.StatusConflict,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal the book to JSON
			bookJSON, err := json.Marshal(tc.book)
			if err != nil {
				t.Fatal(err)
			}
			
			// Create a request with a book in the body
			req, err := http.NewRequest("POST", "/books", bytes.NewBuffer(bookJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handler.CreateBook(rr, req)
			
			// Check status code
			assert.Equal(t, tc.wantStatus, rr.Code)
			
			// For successful requests, verify the book was created
			if tc.wantStatus == http.StatusCreated {
				var book models.Book
				err = json.Unmarshal(rr.Body.Bytes(), &book)
				assert.NoError(t, err)
				assert.Equal(t, tc.book.Title, book.Title)
				
				// Verify the book was added to the store
				createdBook, err := store.GetBook(book.BookID)
				assert.NoError(t, err)
				assert.Equal(t, tc.book.Title, createdBook.Title)
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	// Setup
	store := setupTestBookStore()
	handler := handlers.NewBookHandler(store)
	
	// Test cases
	tests := []struct {
		name       string
		bookID     string
		updateData models.Book
		wantStatus int
	}{
		{
			"Valid Update",
			"1",
			models.Book{
				Title:           "Updated Book 1",
				AuthorID:        "auth-005",
				PublisherID:     "pub-005",
				PublicationDate: "2023-05-15",
				Description:     "Updated description for book 1",
				Price:           24.99,
				ISBN:            "1234567890",
				Pages:           350,
				Genre:           "Sci-Fi",
				Quantity:        8,
			},
			http.StatusOK,
		},
		{
			"Book Not Found",
			"999",
			models.Book{
				Title: "Non-existent Book",
			},
			http.StatusNotFound,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal the update data to JSON
			updateJSON, err := json.Marshal(tc.updateData)
			if err != nil {
				t.Fatal(err)
			}
			
			// Create a request
			req, err := http.NewRequest("PUT", "/books/"+tc.bookID, bytes.NewBuffer(updateJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			
			// Add URL parameters to the request
			vars := map[string]string{
				"id": tc.bookID,
			}
			req = mux.SetURLVars(req, vars)
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handler.UpdateBook(rr, req)
			
			// Check status code
			assert.Equal(t, tc.wantStatus, rr.Code)
			
			// For successful updates, verify the book was updated
			if tc.wantStatus == http.StatusOK {
				var book models.Book
				err = json.Unmarshal(rr.Body.Bytes(), &book)
				assert.NoError(t, err)
				assert.Equal(t, tc.updateData.Title, book.Title)
				assert.Equal(t, tc.bookID, book.BookID) // ID should not change
				
				// Verify the book was updated in the store
				updatedBook, err := store.GetBook(tc.bookID)
				assert.NoError(t, err)
				assert.Equal(t, tc.updateData.Title, updatedBook.Title)
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	// Test cases
	tests := []struct {
		name       string
		bookID     string
		wantStatus int
	}{
		{"Valid ID", "1", http.StatusNoContent},
		{"Invalid ID", "999", http.StatusNotFound},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup a fresh store for each test case
			store := setupTestBookStore()
			handler := handlers.NewBookHandler(store)
			
			// Create a request
			req, err := http.NewRequest("DELETE", "/books/"+tc.bookID, nil)
			if err != nil {
				t.Fatal(err)
			}
			
			// Add URL parameters to the request
			vars := map[string]string{
				"id": tc.bookID,
			}
			req = mux.SetURLVars(req, vars)
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handler.DeleteBook(rr, req)
			
			// Check status code
			assert.Equal(t, tc.wantStatus, rr.Code)
			
			// For successful deletions, verify the book was removed
			if tc.wantStatus == http.StatusNoContent {
				// Try to get the deleted book
				_, err := store.GetBook(tc.bookID)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "book not found")
			}
		})
	}
}