package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/olivier-twist/kindle/internal/reader"
)

func main() {

	// Load the list of books from the file system.
	rootDir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Fatalf("Failed to load books: %s", err)
	}
	filepath := filepath.Join(rootDir, "data", "book_list.txt")

	books, err := reader.ReadBooksFromFileSystem(filepath)
	if err != nil {
		log.Fatalf("Failed to load books: %s", err)
	}

	for _, book := range books {
		fmt.Printf("Title: %s  Author: %s\n\n", book.Title, book.Author)
	}
}
