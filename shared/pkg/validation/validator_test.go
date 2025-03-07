package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	// Valid emails
	validEmails := []string{
		"test@example.com",
		"user.name@subdomain.example.co.uk",
		"user+tag@example.com",
	}
	
	for _, email := range validEmails {
		err := ValidateEmail(email)
		assert.NoError(t, err, "Expected valid email: %s", email)
	}
	
	// Invalid emails
	invalidEmails := []string{
		"",                  // Empty
		"notanemail",        // No @
		"@example.com",      // No local part
		"user@",             // No domain
		"user@invalid",      // No TLD
		"us er@example.com", // Space in local part
		"a@b.c",            // Too short domain parts
	}
	
	for _, email := range invalidEmails {
		err := ValidateEmail(email)
		assert.Error(t, err, "Expected invalid email: %s", email)
	}
}

func TestValidateUsername(t *testing.T) {
	// Valid usernames
	validUsernames := []string{
		"user123",
		"john_doe",
		"alice_789",
		"abc123",
	}
	
	for _, username := range validUsernames {
		err := ValidateUsername(username)
		assert.NoError(t, err, "Expected valid username: %s", username)
	}
	
	// Invalid usernames
	invalidUsernames := []string{
		"",                // Empty
		"ab",              // Too short
		"user-123",        // Contains hyphen
		"user.name",       // Contains period
		"very_long_username_exceeding_thirty_chars", // Too long
	}
	
	for _, username := range invalidUsernames {
		err := ValidateUsername(username)
		assert.Error(t, err, "Expected invalid username: %s", username)
	}
}

func TestValidateName(t *testing.T) {
	// Valid names
	validNames := []string{
		"John",
		"Mary Smith",
		"Jean-Claude",
		"O'Connor",
	}
	
	for _, name := range validNames {
		err := ValidateName(name, "name")
		assert.NoError(t, err, "Expected valid name: %s", name)
	}
	
	// Invalid names
	invalidNames := []string{
		"",           // Empty
		"A",          // Too short
		"John123",    // Contains numbers
		"User@Name",  // Contains special character
	}
	
	for _, name := range invalidNames {
		err := ValidateName(name, "name")
		assert.Error(t, err, "Expected invalid name: %s", name)
	}
}

func TestValidateRequired(t *testing.T) {
	// Valid non-empty strings
	validStrings := []string{
		"Hello",
		"123",
		" trimmed ", // Will be trimmed but not empty
	}
	
	for _, str := range validStrings {
		err := ValidateRequired(str, "field")
		assert.NoError(t, err, "Expected valid non-empty string: %s", str)
	}
	
	// Invalid empty strings
	invalidStrings := []string{
		"",    // Empty
		" ",   // Just whitespace
		"\t",  // Tab
		"\n",  // Newline
	}
	
	for _, str := range invalidStrings {
		err := ValidateRequired(str, "field")
		assert.Error(t, err, "Expected error for empty string: %s", str)
	}
}