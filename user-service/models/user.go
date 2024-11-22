// models/user.go
package models

import (
	"time"

	"github.com/CourtIQ/courtiq-backend/user-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Location represents the database model for location
type Location struct {
	City      string  `bson:"city,omitempty"`
	State     string  `bson:"state,omitempty"`
	Country   string  `bson:"country,omitempty"`
	Latitude  float64 `bson:"latitude,omitempty"`
	Longitude float64 `bson:"longitude,omitempty"`
}

// ToGraphQL converts the DB location to a GraphQL location
func (l *Location) ToGraphQL() *model.Location {
	if l == nil {
		return nil
	}
	return &model.Location{
		City:      &l.City,
		State:     &l.State,
		Country:   &l.Country,
		Latitude:  &l.Latitude,
		Longitude: &l.Longitude,
	}
}

// FromGraphQLInput updates location from GraphQL input
func (l *Location) FromGraphQLInput(input *model.LocationInput) {
	if input == nil {
		return
	}
	if input.City != nil {
		l.City = *input.City
	}
	if input.State != nil {
		l.State = *input.State
	}
	if input.Country != nil {
		l.Country = *input.Country
	}
	if input.Latitude != nil {
		l.Latitude = *input.Latitude
	}
	if input.Longitude != nil {
		l.Longitude = *input.Longitude
	}
}

// User represents the database model for user
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Email          string             `bson:"email"`
	FirstName      string             `bson:"firstName,omitempty"`
	LastName       string             `bson:"lastName,omitempty"`
	DisplayName    string             `bson:"displayName,omitempty"`
	ProfilePicture string             `bson:"profilePicture,omitempty"`
	DateOfBirth    time.Time          `bson:"dateOfBirth,omitempty"`
	Location       *Location          `bson:"location,omitempty"`
	Bio            string             `bson:"bio,omitempty"`
	Rating         int                `bson:"rating,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
}

// ToGraphQL converts the DB user to a GraphQL user
func (u *User) ToGraphQL() *model.User {
	if u == nil {
		return nil
	}

	gqlUser := &model.User{
		ID:             u.ID.Hex(),
		Email:          u.Email,
		FirstName:      &u.FirstName,
		LastName:       &u.LastName,
		DisplayName:    &u.DisplayName,
		ProfilePicture: &u.ProfilePicture,
		Bio:            &u.Bio,
		Rating:         &u.Rating,
		CreatedAt:      &u.CreatedAt,
		UpdatedAt:      &u.UpdatedAt,
	}

	// Convert date of birth if it's not zero
	if !u.DateOfBirth.IsZero() {
		gqlUser.DateOfBirth = &u.DateOfBirth
	}

	// Convert location if it exists
	if u.Location != nil {
		gqlUser.Location = u.Location.ToGraphQL()
	}

	return gqlUser
}

// FromGraphQLInput updates the user from a GraphQL input
func (u *User) FromGraphQLInput(input *model.UpdateUserInput) {
	if input == nil {
		return
	}

	if input.FirstName != nil {
		u.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		u.LastName = *input.LastName
	}
	if input.DisplayName != nil {
		u.DisplayName = *input.DisplayName
	}
	if input.ProfilePicture != nil {
		u.ProfilePicture = *input.ProfilePicture
	}
	if input.Bio != nil {
		u.Bio = *input.Bio
	}
	if input.DateOfBirth != nil {
		u.DateOfBirth = *input.DateOfBirth
	}
	if input.Rating != nil {
		u.Rating = *input.Rating
	}

	// Handle location update
	if input.Location != nil {
		if u.Location == nil {
			u.Location = &Location{}
		}
		u.Location.FromGraphQLInput(input.Location)
	}
}

// NewUser creates a new User instance
func NewUser() *User {
	return &User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewUserFromGraphQLInput creates a new User from GraphQL input
func NewUserFromGraphQLInput(input *model.UpdateUserInput) *User {
	user := NewUser()
	user.FromGraphQLInput(input)
	return user
}
