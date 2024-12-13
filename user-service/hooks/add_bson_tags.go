package main

import (
	"log"
	"os"
	"regexp"
)

func main() {
	filePath := "graph/model/models_gen.go" // Replace with your actual models file path

	// Read the file
	content, err := os.ReadFile(filePath) // Use os.ReadFile instead of ioutil.ReadFile
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Regex to add bson:"_id" to the id field
	re := regexp.MustCompile(`json:"id"`)
	updatedContent := re.ReplaceAll(content, []byte(`json:"id" bson:"_id"`))

	// Write the updated file back
	if err := os.WriteFile(filePath, updatedContent, os.ModePerm); err != nil { // Use os.WriteFile instead of ioutil.WriteFile
		log.Fatalf("Failed to write file: %v", err)
	}

	log.Println("Successfully updated bson tags!")
}
