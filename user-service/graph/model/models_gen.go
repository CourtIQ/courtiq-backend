// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Location struct {
	// Name of the city.
	City *string `json:"city,omitempty"`
	// Name of the state or province.
	State *string `json:"state,omitempty"`
	// Name of the country.
	Country *string `json:"country,omitempty"`
	// Geographical latitude coordinate.
	Latitude *float64 `json:"latitude,omitempty"`
	// Geographical longitude coordinate.
	Longitude *float64 `json:"longitude,omitempty"`
}

type LocationInput struct {
	// Name of the city.
	City *string `json:"city,omitempty"`
	// Name of the state or province.
	State *string `json:"state,omitempty"`
	// Name of the country.
	Country *string `json:"country,omitempty"`
	// Geographical latitude coordinate.
	Latitude *float64 `json:"latitude,omitempty"`
	// Geographical longitude coordinate.
	Longitude *float64 `json:"longitude,omitempty"`
}

type Mutation struct {
}

type Query struct {
}

type UpdateUserInput struct {
	// User's first name.
	FirstName *string `json:"firstName,omitempty"`
	// User's last name.
	LastName *string `json:"lastName,omitempty"`
	// User's chosen display name or username.
	DisplayName *string `json:"displayName,omitempty"`
	// URL to the user's profile picture.
	ProfilePicture *string `json:"profilePicture,omitempty"`
	// User's date of birth.
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty"`
	// Short biography or description provided by the user.
	Bio *string `json:"bio,omitempty"`
	// User's rating within the app.
	Rating *int `json:"rating,omitempty"`
	// User's geographical location.
	Location *LocationInput `json:"location,omitempty"`
}

type User struct {
	// Unique identifier for the user.
	ID string `json:"id"`
	// User's email address. Accessible only by the user themselves.
	Email string `json:"email"`
	// User's first name.
	FirstName *string `json:"firstName,omitempty"`
	// User's last name.
	LastName *string `json:"lastName,omitempty"`
	// User's chosen display name or username.
	DisplayName *string `json:"displayName,omitempty"`
	// URL to the user's profile picture.
	ProfilePicture *string `json:"profilePicture,omitempty"`
	// User's date of birth.
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty"`
	// User's geographical location.
	Location *Location `json:"location,omitempty"`
	// Short biography or description provided by the user.
	Bio *string `json:"bio,omitempty"`
	// User's rating within the app.
	Rating *int `json:"rating,omitempty"`
	// Timestamp when the user account was created.
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// Timestamp when the user account was last updated.
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (User) IsEntity() {}
