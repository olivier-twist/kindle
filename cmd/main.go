package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/olivier-twist/kindle/internal/openapi"
	"github.com/olivier-twist/kindle/internal/reader"
)

// GetFilePath returns a path to an input file.
func GetFilePath(filename string) (string, error) {
	rootDir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return "", fmt.Errorf("Unable to get rootdir")
	}
	filepath := filepath.Join(rootDir, "data", filename)
	return filepath, nil
}

/*
func main() {
	apiKey := os.Getenv("OPENAPI_KEY")         // Replace with your OpenAI API key
	filePath, err := GetFilePath("book.jsonl") // Replace with the path to your file
	purpose := "user_data"                     // Replace with the purpose of the file, e.g., "fine-tune"

	if err != nil {
		log.Fatalf("%v", err)
	}
	err = openapi.UploadFile(apiKey, filePath, purpose)
	if err != nil {
		fmt.Printf("Error uploading file: %v\n", err)
	} else {
		fmt.Println("File uploaded successfully.")
	}
}
*/

func main() {
	apiKey := os.Getenv("OPENAPI_KEY")
	bookPath, err := GetFilePath("book.jsonl")

	if err != nil {
		log.Fatalf("%v", err)
	}
	booksJSON, err := reader.GetBooksFromJsonFile(bookPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	tags, err := GetFilePath("tag.jsonl")
	tagsJSON, err := reader.GetTagsFromJsonFile(tags)
	if err != nil {
		log.Fatalf("%v", err)
	}

	res, err := openapi.AssignTagsToBooks(apiKey, booksJSON, tagsJSON)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf(res)
}
