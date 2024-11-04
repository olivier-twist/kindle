package dbops

import (
	"database/sql"
	"fmt"

	"github.com/olivier-twist/kindle/internal/common"
)

func InsertBooks(db *sql.DB, books []common.Book) error {

	stmt, err := db.Prepare("Insert into BOOK (AUTHOR, TITLE) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("unable to prepare statement %v", err)
	}

	defer stmt.Close()

	for _, book := range books {

		_, err = stmt.Exec(book.Author, book.Title)
		if err != nil {
			return fmt.Errorf("failed to insert book: %v %v", book, err)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to insert book %v", err)
	}

	return nil
}

func InsertTags(db *sql.DB, tags []string) error {
	stmt, err := db.Prepare("INSERT INTO TAG (TAG) VALUES (?)")
	if err != nil {
		return fmt.Errorf("failed to create Prepared Statement for tag insertion : %v", err)
	}

	defer stmt.Close()

	for _, tag := range tags {
		_, err = stmt.Exec(tag)

		if err != nil {
			return fmt.Errorf("failed to insert tag %v error - %v", tag, err)
		}
	}

	return nil
}

func InsertBookTags(db *sql.DB, book_tag map[string][]string) error {
	stmt, err := db.Prepare("INSERT INTO BOOK_TAG (BOOK, TAG) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to create Prepared Statement for book tag insertion : %v", err)
	}
	defer stmt.Close()

	for book, tags := range book_tag {
		for _, tag := range tags {
			_, err = stmt.Exec(book, tag)
			if err != nil {
				return fmt.Errorf("failed to insert book tag %v - %v error - %v", book, tag, err)
			}
		}
	}
	return nil

}
