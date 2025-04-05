package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/foverokavindz/book-api/models"
)

// BookSearchResult represents a book with its search relevance score
type BookSearchResult struct {
	Book  models.Book
	Score int
}

// SearchBooks handles GET /books/search?q={keyword}
func (h *BookHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {

	query := strings.ToLower(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	books, err := h.store.GetAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Use concurrent search for optimization
	results := ConcurrentSearch(books, query)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// ConcurrentSearch performs on books based on title and description
func ConcurrentSearch(books []models.Book, query string) []models.Book {
	// Create channels for communication between goroutines
	resultsChannel := make(chan BookSearchResult)
	var wg sync.WaitGroup

	// Launch a goroutine for each book for parallel processing
	for _, book := range books {
		wg.Add(1)
		go func(b models.Book) {
			defer wg.Done()
			score := CalculateSearchScore(b, query)
			if score > 0 {
				resultsChannel <- BookSearchResult{Book: b, Score: score}
			}
		}(book)
	}

	// Close the results channel once all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChannel)
	}()

	// Collect results from channel
	var searchResults []BookSearchResult
	for result := range resultsChannel {
		searchResults = append(searchResults, result)
	}

	// Sort results by score in descending order
	sort.Slice(searchResults, func(i, j int) bool {
		return searchResults[i].Score > searchResults[j].Score
	})

	// Extract the books from results
	var result []models.Book
	for _, searchResult := range searchResults {
		result = append(result, searchResult.Book)
	}

	return result
}

// CalculateSearchScore calculates a relevance score for a book based on the search query
func CalculateSearchScore(book models.Book, query string) int {
	score := 0
	title := strings.ToLower(book.Title)
	description := strings.ToLower(book.Description)

	// Check if query appears in the title - give more weight to exact matches
	if strings.Contains(title, query) {
		score += 10
	}

	// Check how many words from the query appear in the title
	queryWords := strings.Fields(query)

	for _, word := range queryWords {
		if strings.Contains(title, word) {
			score += 3
		}
	}

	// Check if query appears in the description
	if strings.Contains(description, query) {
		score += 5
	}

	// Check how many words from the query appear in the description
	for _, word := range queryWords {
		if strings.Contains(description, word) {
			score += 1
		}
	}

	return score
}