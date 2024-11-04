package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/olivier-twist/kindle/internal/dbops"
	"github.com/olivier-twist/kindle/internal/openapi"
	"github.com/olivier-twist/kindle/internal/reader"
)

// GetFilePath returns a path to an input file.
func GetFilePath(filename string) (string, error) {
	rootDir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return "", fmt.Errorf("unable to get rootdir")
	}
	filepath := filepath.Join(rootDir, "data", filename)
	return filepath, nil
}

/*
func main() {
	godotenv.Load()
	db_user := os.Getenv("DB_USER")
	db_pwd := os.Getenv("DB_PWD")
	bookPath, err := GetFilePath("book_list.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}

	books, err := reader.ReadBooksFromTxtFile(bookPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	//database driver
	db, err := sql.Open("mysql", db_user+":"+db_pwd+"@tcp(127.0.0.1:3306)/book")

	if err != nil {
		log.Fatalf("**%v", err)
	}

	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	err = dbops.InsertBooks(db, books)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
*/

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load env file")
	}

	apiKey := os.Getenv("OPENAPI_KEY")
	db_user := os.Getenv("DB_USER")
	db_pwd := os.Getenv("DB_PWD")

	//	bookPath, err := GetFilePath("book.jsonl")

	//	if err != nil {
	//log.Fatalf("%v", err)
	//}
	//books, err := reader.GetBooksFromJsonFile(bookPath)
	//if err != nil {
	//log.Fatalf("%v", err)
	//	}

	//database driver
	db, err := sql.Open("mysql", db_user+":"+db_pwd+"@tcp(127.0.0.1:3306)/book")

	if err != nil {
		log.Fatalf("**%v", err)
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	books, err := reader.GetBooksToBeProcessed(db)

	if err != nil {
		log.Fatalf("%v", err)
	}

	increment := 40
	len := len(books)
	bottom := 0
	top := increment

	// Loop through the books in increments of 40
	for {
		// Update the books from 0 to len(books) in increments of 40
		t := min(top, len)
		res, err := openapi.AssignTagsToBooks(apiKey, books[bottom:t])
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = dbops.InsertBookTags(db, res)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if t == len {
			break
		}
		top += increment
		bottom += increment
	}

}
