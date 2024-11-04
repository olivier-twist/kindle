// Package reader loads the list of books from the file system or from the database.
package reader

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/olivier-twist/kindle/internal/common"
)

// ReadBooksFromFileSystem reads the list of books from the file system.
// path is the file path where the book files are located.
// This code is throw away as it is specific to a given file.
func ReadBooksFromTxtFile(path string) ([]common.Book, error) {
	books := make([]common.Book, 0, 600)

	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("invalid empty file path")
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		log.Printf("%s", err)
		return nil, fmt.Errorf("Failed to open file %s", path)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	inTitle := true
	book := common.Book{}
	isPopulated := false

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if inTitle {
			book.Title = line
			inTitle = false
			isPopulated = false
		} else {
			book.Author = line
			inTitle = true
			isPopulated = true
		}

		if isPopulated {
			books = append(books, book)
			book = common.Book{}
		}
	}
	return books, nil

}

func GetBooksFromJsonFile(path string) ([]common.Book, error) {
	file, err := os.ReadFile(path)

	if err != nil {
		log.Printf("Failed to retrieve tag file")
		return nil, err
	}

	var books []common.Book

	err = json.Unmarshal(file, &books)
	if err != nil {
		log.Printf("Failed to unmarshal file %s into an array of tags ", path)
		return nil, err
	}

	return books, nil
}

func GetBooksToBeProcessed(db *sql.DB) ([]common.Book, error) {
	books := make([]common.Book, 0, 600)
	rows, err := db.Query("SELECT TITLE, AUTHOR FROM BOOK WHERE TITLE NOT IN (SELECT DISTINCT(BOOK) FROM BOOK_TAG)")

	if err != nil {
		log.Printf("Failed to retrieve tag file")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var book common.Book
		if err := rows.Scan(&book.Title, &book.Author); err != nil {
			log.Printf("Failed to retrieve tag file")
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}
