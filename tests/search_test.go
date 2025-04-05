package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/foverokavindz/book-api/handlers"
	"github.com/foverokavindz/book-api/models"
	"github.com/foverokavindz/book-api/storage"
	"github.com/stretchr/testify/assert"
)

// setupTestBooksForSearch creates a test store with books specifically for search tests
func setupTestBooksForSearch() *storage.BookStore {
	store := storage.NewBookStore("search-test-storage")
	
	// Add sample books for search testing
	store.AddBook(models.Book{
		BookID:          "1",
		Title:           "Go Programming",
		AuthorID:        "auth-101",
		PublisherID:     "pub-101",
		PublicationDate: "2022-10-15",
		Description:     "Learn Go programming language basics and advanced concepts",
		Price:           29.99,
		ISBN:            "1111111111",
		Pages:           450,
		Genre:           "Programming",
		Quantity:        20,
	})
	
	store.AddBook(models.Book{
		BookID:          "2",
		Title:           "Advanced Python",
		AuthorID:        "auth-102",
		PublisherID:     "pub-102",
		PublicationDate: "2022-11-20",
		Description:     "This book covers Python and some Go concepts",
		Price:           34.99,
		ISBN:            "2222222222",
		Pages:           380,
		Genre:           "Programming",
		Quantity:        15,
	})
	
	store.AddBook(models.Book{
		BookID:          "3",
		Title:           "Database Design",
		AuthorID:        "auth-103",
		PublisherID:     "pub-103",
		PublicationDate: "2022-09-10",
		Description:     "Comprehensive guide to database design and implementation",
		Price:           24.99,
		ISBN:            "3333333333",
		Pages:           320,
		Genre:           "Database",
		Quantity:        12,
	})
	
	return store
}

func TestSearchBooks(t *testing.T) {
	// Setup
	test_store := setupTestBooksForSearch()
	handler := handlers.NewBookHandler(test_store)
	
	// Test cases
	tests := []struct {
		name           string
		query          string
		wantStatus     int
		expectedCount  int
		expectedTitles []string
	}{
		{
			"Search by title exact match",
			"Go Programming",
			http.StatusOK,
			2,
			[]string{"Go Programming"},
		},
		{
			"Search by partial title",
			"Go",
			http.StatusOK,
			2,
			[]string{"Go Programming", "Advanced Python"}, // "Go Programming" should come first due to higher relevance
		},
		{
			"Search by description",
			"database",
			http.StatusOK,
			1,
			[]string{"Database Design"},
		},
		{
			"Search with no results",
			"javascript",
			http.StatusOK,
			0,
			[]string{},
		},
		{
			"Empty search query",
			"",
			http.StatusBadRequest,
			0,
			[]string{},
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with query parameter
			req, err := http.NewRequest("GET", "/books/search?q="+tc.query, nil)
			if err != nil {
				t.Fatal(err)
			}
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Call the handler
			handler.SearchBooks(rr, req)
			
			// Check status code
			assert.Equal(t, tc.wantStatus, rr.Code)
			
			// For successful requests, check the search results
			if tc.wantStatus == http.StatusOK {
				var books []models.Book
				err = json.Unmarshal(rr.Body.Bytes(), &books)
				assert.NoError(t, err)
				
				// Check count of results
				assert.Equal(t, tc.expectedCount, len(books))
				
				// Check if expected titles match
				for i := 0; i < len(tc.expectedTitles) && i < len(books); i++ {
					assert.Equal(t, tc.expectedTitles[i], books[i].Title)
				}
			}
		})
	}
}

// Test the concurrentSearch function
func TestConcurrentSearch(t *testing.T) {
	// Set up test books
	books := []models.Book{
		{
			BookID:          "1",
			Title:           "Go Programming",
			Description:     "Learn Go programming language basics",
			PublicationDate: "2022-01-10",
			Pages:           300,
			Genre:           "Programming",
			Quantity:        15,
		},
		{
			BookID:          "2",
			Title:           "Advanced Programming Topics",
			Description:     "This covers some go concepts and other languages",
			PublicationDate: "2022-02-15",
			Pages:           280,
			Genre:           "Programming",
			Quantity:        12,
		},
		{
			BookID:          "3",
			Title:           "Database Design",
			Description:     "Comprehensive guide to database design",
			PublicationDate: "2022-03-20",
			Pages:           320,
			Genre:           "Database",
			Quantity:        10,
		},
	}
	
	// Test cases
	tests := []struct {
		name           string
		query          string
		expectedCount  int
		expectedFirst  string // BookID of expected first result
	}{
		{"Search for Go", "go", 2, "1"}, // Should find books 1 and 2, with 1 having higher relevance
		{"Search for Programming", "programming", 2, "1"}, // Should find books 1 and 2
		{"Search for Database", "database", 1, "3"}, // Should find only book 3
		{"No results", "javascript", 0, ""}, // Should find no books
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results := handlers.ConcurrentSearch(books, tc.query)
			
			// Check count of results
			assert.Equal(t, tc.expectedCount, len(results))
			
			// Check first result if expected
			if tc.expectedCount > 0 && tc.expectedFirst != "" {
				assert.Equal(t, tc.expectedFirst, results[0].BookID)
			}
		})
	}
}