package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/foverokavindz/book-api/models"
	"github.com/foverokavindz/book-api/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BookHandler struct {
	store *storage.BookStore
}

// NewBookHandler creates a new BookHandler with the given BookStore
func NewBookHandler(store *storage.BookStore) *BookHandler {
	return &BookHandler{store: store}
}

// GetBooks handles GET /books
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {

	// Parse pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Default values
	limit := 2
	offset := 0


    // Parse limit
    if limitStr != "" {
        parsedLimit, err := strconv.Atoi(limitStr)
        if err != nil || parsedLimit < 1 {
            http.Error(w, "Invalid limit parameter, must be a positive integer", http.StatusBadRequest)
            return
        }
        limit = parsedLimit
    }
    
    // Parse offset
    if offsetStr != "" {
        parsedOffset, err := strconv.Atoi(offsetStr)
        if err != nil || parsedOffset < 0 {
            http.Error(w, "Invalid offset parameter, must be a non-negative integer", http.StatusBadRequest)
            return
        }
        offset = parsedOffset
    }
	

	books, err := h.store.GetAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	// Calculate total books and pages
	totalBooks := len(books)
	totalPages := (totalBooks + limit - 1) / limit 


	// Apply pagination
	end := offset + limit
	if end > totalBooks {
		end = totalBooks
	}
	
	var paginatedBooks []models.Book
	if offset < totalBooks {
		paginatedBooks = books[offset:end]
	} else {
		paginatedBooks = []models.Book{} // Empty slice if offset exceeds total
	}
			

	// Create response with metadata
	response := struct {
		Books      []models.Book `json:"books"`
		Pagination struct {
			Total       int `json:"total"`
			Limit       int `json:"limit"`
			Offset      int `json:"offset"`
			TotalPages  int `json:"totalPages"`
			CurrentPage int `json:"currentPage"`
		} `json:"pagination"`
	}{
		Books: paginatedBooks,
	}

	// Fill pagination metadata
	response.Pagination.Total = totalBooks
	response.Pagination.Limit = limit
	response.Pagination.Offset = offset
	response.Pagination.TotalPages = totalPages
	response.Pagination.CurrentPage = (offset / limit) + 1

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// GetBook handles GET /books/{id}
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	// Extract the book ID from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	book, err := h.store.GetBook(id)
	if err != nil {
		if err.Error() == "book not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// CreateBook handles POST /books
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	// Decode the request body into the Book struct
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Generate a new UUID if not provided
	if book.BookID == "" {
		book.BookID = uuid.New().String()
	}

	err = h.store.AddBook(book)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// UpdateBook handles PUT /books/{id}
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// prevent user from changing the book ID
	book.BookID = id

	err = h.store.UpdateBook(id, book)
	if err != nil {
		if err.Error() == "book not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// DeleteBook handles DELETE /books/{id}
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.store.DeleteBook(id)
	if err != nil {
		if err.Error() == "book not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}