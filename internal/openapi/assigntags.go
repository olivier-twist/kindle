// Package openapi assigns relevant tags to a given book
package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/olivier-twist/kindle/internal/common"
)

const (
	uploadURL         = "https://api.openai.com/v1/files"
	chatCompletionURL = "https://api.openai.com/v1/chat/completions"
)

var openapi_key = os.Getenv("OPENAPI_KEY")

// UploadFile uploads a file to OpenApi
func UploadFile(apiKey, filePath, purpose string) error {
	// Open the file to upload
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Create a buffer to hold the form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file to the form
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("could not create form file: %v", err)
	}

	// Copy the file content to the form file
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("could not copy file content: %v", err)
	}

	// Add the purpose field (required by OpenAI)
	err = writer.WriteField("purpose", purpose)
	if err != nil {
		return fmt.Errorf("could not write purpose field: %v", err)
	}

	// Close the writer to finalize the form data
	writer.Close()

	// Prepare the request
	url := "https://api.openai.com/v1/files"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("could not create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload file, status: %v, response: %s", resp.Status, respBody)
	}

	// Print success message
	fmt.Println("File uploaded successfully")
	return nil
}

// preparePrompt generates a prompt string for the API request
func preparePrompt(books []common.Book, tags []common.Tag) string {
	sort.Sort(common.ById(books))
	// Convert books and tags to JSON format for prompt
	booksJSON, _ := json.Marshal(books)
	tagsJSON, _ := json.Marshal(tags)

	// Create the prompt
	prompt := fmt.Sprintf(`Given the following array of books and array of tags:
						Books:
						%s

						Tags:
						%s

						Given a list of books, each defined by an id, title, and one or more authors, find relevant information about each book online using OpenAPI. Analyze the content, themes, genre, and topics associated with each book. Based on this analysis, return a map where each book title is a key, and the corresponding value is a set of tag  that best describe the book’s content and themes. Add a fiction tag if the genre is fiction. Each tag id corresponds to a predefined list of tags in the format {id: int, tag: string}`, booksJSON, tagsJSON)

	return prompt
}

// Assign Tags to Books
func AssignTagsToBooks(apiKey string, books []common.Book, tags []common.Tag) (string, error) {
	// Prepare the prompt with book and tag data
	prompt := preparePrompt(books, tags)

	// Set up the request payload
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo", // Replace with your preferred model
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("could not create request body: %v", err)
	}

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status and read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Print the error message from the API
		return "", fmt.Errorf("failed to get response, status: %v, response: %s", resp.Status, string(body))
	}

	// Parse the response JSON
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("could not parse response: %v", err)
	}

	// Check if "choices" is in the response
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("response does not contain choices: %v", response)
	}

	// Extract the message content from the first choice
	content := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return content, nil
}
