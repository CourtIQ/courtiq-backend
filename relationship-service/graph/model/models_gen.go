// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Relationship interface defines the common structure for all relationship types
// between two users (Friendship, Coachship, etc.).
type Relationship interface {
	IsRelationship()
	// The unique identifier for the relationship.
	GetID() primitive.ObjectID
	// The type of the relationship (e.g., FRIENDSHIP, COACHSHIP).
	GetType() RelationshipType
	// The current status of the relationship (e.g., PENDING, ACCEPTED, BLOCKED).
	GetStatus() RelationshipStatus
	// The user who initiated the relationship or the request.
	GetInitiator() *User
	// The user who received the relationship request or is the target.
	GetReceiver() *User
	// Timestamp when the relationship was first created or requested.
	GetCreatedAt() time.Time
	// Timestamp when the relationship was last updated (e.g., status change).
	GetUpdatedAt() *time.Time
}

// Represents a coaching relationship between two users (a coach and a student).
type Coachship struct {
	// The unique identifier for the coachship.
	ID primitive.ObjectID `json:"id" bson:"_id"`
	// Always COACHSHIP for this type.
	Type RelationshipType `json:"type" bson:"type"`
	// The current status of the coachship (e.g., PENDING, ACCEPTED).
	Status RelationshipStatus `json:"status" bson:"status"`
	// The user who initiated the coaching request.
	Initiator *User `json:"initiator" bson:"initiator"`
	// The user who received the coaching request.
	Receiver *User `json:"receiver" bson:"receiver"`
	// Timestamp when the coaching request was sent or created.
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	// Timestamp when the coachship status last changed.
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	// The user acting as the coach in this relationship.
	Coach *User `json:"coach" bson:"coach"`
	// The user acting as the student in this relationship.
	Student *User `json:"student" bson:"student"`
}

func (Coachship) IsRelationship() {}

// The unique identifier for the relationship.
func (this Coachship) GetID() primitive.ObjectID { return this.ID }

// The type of the relationship (e.g., FRIENDSHIP, COACHSHIP).
func (this Coachship) GetType() RelationshipType { return this.Type }

// The current status of the relationship (e.g., PENDING, ACCEPTED, BLOCKED).
func (this Coachship) GetStatus() RelationshipStatus { return this.Status }

// The user who initiated the relationship or the request.
func (this Coachship) GetInitiator() *User { return this.Initiator }

// The user who received the relationship request or is the target.
func (this Coachship) GetReceiver() *User { return this.Receiver }

// Timestamp when the relationship was first created or requested.
func (this Coachship) GetCreatedAt() time.Time { return this.CreatedAt }

// Timestamp when the relationship was last updated (e.g., status change).
func (this Coachship) GetUpdatedAt() *time.Time { return this.UpdatedAt }

func (Coachship) IsEntity() {}

// Represents a friendship relationship between two users.
type Friendship struct {
	// The unique identifier for the friendship.
	ID primitive.ObjectID `json:"id" bson:"_id"`
	// Always FRIENDSHIP for this type.
	Type RelationshipType `json:"type" bson:"type"`
	// The current status of the friendship (e.g., PENDING, ACCEPTED, BLOCKED).
	Status RelationshipStatus `json:"status" bson:"status"`
	// The user who sent the friend request.
	Initiator *User `json:"initiator" bson:"initiator"`
	// The user who received the friend request.
	Receiver *User `json:"receiver" bson:"receiver"`
	// Timestamp when the friendship request was sent or created.
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	// Timestamp when the friendship status last changed.
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func (Friendship) IsRelationship() {}

// The unique identifier for the relationship.
func (this Friendship) GetID() primitive.ObjectID { return this.ID }

// The type of the relationship (e.g., FRIENDSHIP, COACHSHIP).
func (this Friendship) GetType() RelationshipType { return this.Type }

// The current status of the relationship (e.g., PENDING, ACCEPTED, BLOCKED).
func (this Friendship) GetStatus() RelationshipStatus { return this.Status }

// The user who initiated the relationship or the request.
func (this Friendship) GetInitiator() *User { return this.Initiator }

// The user who received the relationship request or is the target.
func (this Friendship) GetReceiver() *User { return this.Receiver }

// Timestamp when the relationship was first created or requested.
func (this Friendship) GetCreatedAt() time.Time { return this.CreatedAt }

// Timestamp when the relationship was last updated (e.g., status change).
func (this Friendship) GetUpdatedAt() *time.Time { return this.UpdatedAt }

func (Friendship) IsEntity() {}

// Provides structured geographical details about a user's location.
// All fields are optional and can be omitted if unknown.
type Location struct {
	City      *string  `json:"city,omitempty" bson:"city,omitempty"`
	State     *string  `json:"state,omitempty" bson:"state,omitempty"`
	Country   *string  `json:"country,omitempty" bson:"country,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
}

type Mutation struct {
}

type Query struct {
}

type RelationshipQueries struct {
	GetFriendship *Friendship `json:"getFriendship,omitempty" bson:"getFriendship,omitempty"`
}

type User struct {
	// The unique identifier for the user.
	ID primitive.ObjectID `json:"id" bson:"_id"`
	// User's first name.
	FirstName *string `json:"firstName,omitempty" bson:"firstName,omitempty"`
	// User's last name.
	LastName *string `json:"lastName,omitempty" bson:"lastName,omitempty"`
	// User's display name.
	DisplayName *string `json:"displayName,omitempty" bson:"displayName,omitempty"`
	// User's chosen username.
	Username *string `json:"username,omitempty" bson:"username,omitempty"`
	// URL to the user's profile picture.
	ProfilePicture *string `json:"profilePicture,omitempty" bson:"profilePicture,omitempty"`
}

func (User) IsEntity() {}

type RelationshipStatus string

const (
	RelationshipStatusNone     RelationshipStatus = "NONE"
	RelationshipStatusPending  RelationshipStatus = "PENDING"
	RelationshipStatusAccepted RelationshipStatus = "ACCEPTED"
	RelationshipStatusDeclined RelationshipStatus = "DECLINED"
	RelationshipStatusBlocked  RelationshipStatus = "BLOCKED"
)

var AllRelationshipStatus = []RelationshipStatus{
	RelationshipStatusNone,
	RelationshipStatusPending,
	RelationshipStatusAccepted,
	RelationshipStatusDeclined,
	RelationshipStatusBlocked,
}

func (e RelationshipStatus) IsValid() bool {
	switch e {
	case RelationshipStatusNone, RelationshipStatusPending, RelationshipStatusAccepted, RelationshipStatusDeclined, RelationshipStatusBlocked:
		return true
	}
	return false
}

func (e RelationshipStatus) String() string {
	return string(e)
}

func (e *RelationshipStatus) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RelationshipStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RelationshipStatus", str)
	}
	return nil
}

func (e RelationshipStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RelationshipType string

const (
	RelationshipTypeFriendship RelationshipType = "FRIENDSHIP"
	RelationshipTypeCoachship  RelationshipType = "COACHSHIP"
)

var AllRelationshipType = []RelationshipType{
	RelationshipTypeFriendship,
	RelationshipTypeCoachship,
}

func (e RelationshipType) IsValid() bool {
	switch e {
	case RelationshipTypeFriendship, RelationshipTypeCoachship:
		return true
	}
	return false
}

func (e RelationshipType) String() string {
	return string(e)
}

func (e *RelationshipType) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RelationshipType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RelationshipType", str)
	}
	return nil
}

func (e RelationshipType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
