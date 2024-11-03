package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/olivier-twist/kindle/internal/common"
)

func ReadTagsFromFile(path string) ([]string, error) {
	tags := make([]string, 0, 0)

	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("No path specified")
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file %v", path)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		tag := strings.TrimSpace(scanner.Text())
		tags = append(tags, tag)
	}

	return tags, nil

}

func GetTagsFromJsonFile(path string) ([]common.Tag, error) {
	file, err := os.ReadFile(path)

	if err != nil {
		log.Printf("Failed to retrieve tag file")
		return nil, err
	}

	var tags []common.Tag

	err = json.Unmarshal(file, &tags)
	if err != nil {
		log.Printf("Failed to unmarshal file %s into an array of tags ", path)
		return nil, err
	}

	return tags, nil
}
