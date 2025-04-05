package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/foverokavindz/book-api/handlers"
	"github.com/foverokavindz/book-api/storage"
	"github.com/gorilla/mux"
)

func main() {

	// Create data folder if it doesn't exist
	dataDir := "data"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err := os.Mkdir(dataDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}

	// Initialize book store
	bookStore := storage.NewBookStore(filepath.Join(dataDir, "books.json"))
	// Make storage reference available to handlers
	bookHandler := handlers.NewBookHandler(bookStore)

	// Create router and register routes
	r := mux.NewRouter()
	r.HandleFunc("/books", bookHandler.GetBooks).Methods("GET")
	r.HandleFunc("/books", bookHandler.CreateBook).Methods("POST")
	r.HandleFunc("/books/search", bookHandler.SearchBooks).Methods("GET")
	r.HandleFunc("/books/{id}", bookHandler.GetBook).Methods("GET")
	r.HandleFunc("/books/{id}", bookHandler.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", bookHandler.DeleteBook).Methods("DELETE")

	// Start server
	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}