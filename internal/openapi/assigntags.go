// Package openapi assigns relevant tags to a given book
package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/olivier-twist/kindle/internal/common"
)

// preparePrompt generates a prompt string for the API request
func preparePrompt(books []common.Book) string {
	sort.Sort(common.ById(books))
	// Convert books and tags to JSON format for prompt
	booksJSON, _ := json.Marshal(books)

	// Create the prompt
	prompt := fmt.Sprintf(`Given the following array of books and array of tags:
						Books:
						%s

						Given a list of books, each defined by an id, title, and one or more authors, find relevant information about each book online. Analyze the content, themes, genre, and topics associated with each book. Return output is map[string][]string. The map key is the book title and the value is the set of strinngs  that best describe the book’s content and themes.`, booksJSON)

	return prompt
}

// Assign Tags to Books
func AssignTagsToBooks(apiKey string, books []common.Book) (map[string][]string, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is empty")
	}

	// Prepare the prompt with book data
	prompt := preparePrompt(books)

	// Set up the request payload
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo", // Replace with your preferred model
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("could not create request body: %v", err)
	}

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status and read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Print the error message from the API
		return nil, fmt.Errorf("failed to get response, status: %v, response: %s", resp.Status, string(body))
	}

	// Parse the response JSON
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("assignTagsToBook: could not unmarshal response: %v", err)
	}

	// Check if "choices" is in the response
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("response does not contain choices: %v", response)
	}

	// Extract the message content from the first choice
	content := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return parseBookTagResponse(content)
}

// Parse the response String and return a map of book titles to tags
func parseBookTagResponse(response string) (map[string][]string, error) {
	// Parse the response JSON
	var responseJSON map[string]interface{}
	var m map[string][]string = make(map[string][]string)

	if err := json.Unmarshal([]byte(response), &responseJSON); err != nil {
		fmt.Printf("error: %v\n\n %v", err, response)
		return m, nil /*fmt.Errorf("parseBookTag: could not unmarshal response:%v %v\n\n %v",
		responseJSON, err, response)*/
	}

	for k, v := range responseJSON {
		vinterface, ok := v.([]interface{})
		if !ok {
			fmt.Printf("v is not a []interface{}: %v %v\n\n", v, reflect.TypeOf(v))
			continue
		}
		if len(vinterface) == 0 {
			continue
		}

		for _, v1 := range v.([]interface{}) {
			val, ok := v1.(string)
			if !ok {
				fmt.Printf("v1 is not a string: %v %v\n\n", v1, reflect.TypeOf(v1))
				return m, fmt.Errorf("parseBookTag: could not unmarshal response:%v ",
					responseJSON)

				//continue
			}
			m[k] = append(m[k], strings.ToLower(val))
		}
	}
	return m, nil

}
