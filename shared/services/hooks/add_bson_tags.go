package main

import (
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	filePath := "graph/model/models_gen.go"

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Convert content to string for easier manipulation
	strContent := string(content)

	// Find all json tags
	re := regexp.MustCompile(`json:"([^"]+?)(?:,omitempty)?"`)

	// Replace each json tag with json+bson tags
	strContent = re.ReplaceAllStringFunc(strContent, func(match string) string {
		// Extract the field name and omitempty flag
		fieldMatch := re.FindStringSubmatch(match)
		if len(fieldMatch) < 2 {
			return match
		}

		fieldName := fieldMatch[1]
		hasOmitempty := strings.Contains(match, ",omitempty")

		// Special case for "id" field
		if fieldName == "id" {
			return `json:"id" bson:"_id"`
		}

		// Build the new tag string for other fields
		if hasOmitempty {
			return `json:"` + fieldName + `,omitempty" bson:"` + fieldName + `,omitempty"`
		}
		return `json:"` + fieldName + `" bson:"` + fieldName + `"`
	})

	// Write the updated content back to file
	if err := os.WriteFile(filePath, []byte(strContent), os.ModePerm); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	log.Println("Successfully updated bson tags!")
}
