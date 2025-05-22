package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/foverokavindz/book-api/models"
)

// provide a thread-safe file-based storage for books
type BookStore struct {
	mutex    sync.RWMutex
	filePath string
}
// NewBookStore creates a new BookStore instance
func NewBookStore(filePath string) *BookStore {
	return &BookStore{
		filePath: filePath,
	}
}

// LoadBooks loads all books from the JSON file
func (s *BookStore) LoadBooks() ([]models.Book, error) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()

    // Create empty books.json file if it doesn't exist
    if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
        if fileCreationError := os.WriteFile(s.filePath, []byte("[]"), 0644); fileCreationError != nil {
            return nil, fmt.Errorf("failed to create books file: %w", fileCreationError)
        }
    }

    // Read and parse file in one step when possible
    var books []models.Book

    booksData, conversionError := os.ReadFile(s.filePath)

	// if err is nil, it means the file was read successfully
    if conversionError != nil {
        return nil, fmt.Errorf("failed to read books file: %w", conversionError)
    }
    
    if err := json.Unmarshal(booksData, &books); err != nil {
        return nil, fmt.Errorf("invalid JSON in books file: %w", err)
    }

    return books, nil
}

// SaveBooks saves books to the JSON file
func (s *BookStore) SaveBooks(books []models.Book) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Convert books to formatted JSON and write to file in one logical operation
	jsonData, conversionError := json.MarshalIndent(books, "", "  ")
	if conversionError != nil {
		return fmt.Errorf("failed to convert books to JSON: %w", conversionError)
	}

	if writeError := os.WriteFile(s.filePath, jsonData, 0644); writeError != nil {
		return fmt.Errorf("failed to write books to file: %w", writeError)
	}
	
	return nil
}

// GetAllBooks returns all books
func (s *BookStore) GetAllBooks() ([]models.Book, error) {
	return s.LoadBooks()
}

// GetBook returns a single book by ID
func (s *BookStore) GetBook(id string) (models.Book, error) {
	books, err := s.LoadBooks()
	if err != nil {
		return models.Book{}, err
	}

	for _, book := range books {
		if book.BookID == id {
			return book, nil
		}
	}

	return models.Book{}, errors.New("book not found")
}

// AddBook adds a new book
func (s *BookStore) AddBook(book models.Book) error {
	books, err := s.LoadBooks()
	if err != nil {
		return err
	}

	// Check for duplicate book ID
	for _, b := range books {
		if b.BookID == book.BookID {
			return errors.New("book with this ID already exists")
		}
	}

	books = append(books, book)
	return s.SaveBooks(books)
}

// UpdateBook updates an existing book
func (s *BookStore) UpdateBook(id string, updatedBook models.Book) error {
	books, err := s.LoadBooks()
	if err != nil {
		return err
	}

	found := false
	for i, book := range books {
		if book.BookID == id {
			books[i] = updatedBook
			found = true
			break
		}
	}

	if !found {
		return errors.New("book not found")
	}

	return s.SaveBooks(books)
}

// DeleteBook deletes a book
func (s *BookStore) DeleteBook(id string) error {
	books, err := s.LoadBooks()
	if err != nil {
		return err
	}

	found := false
	var updatedBooks []models.Book
	for _, book := range books {
		if book.BookID != id {
			updatedBooks = append(updatedBooks, book)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("book not found")
	}

	return s.SaveBooks(updatedBooks)
}