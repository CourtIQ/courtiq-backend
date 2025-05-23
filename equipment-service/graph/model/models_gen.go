// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Equipment interface {
	IsEntity()
	IsEquipment()
	GetID() primitive.ObjectID
	GetOwnerID() primitive.ObjectID
	GetName() string
	GetType() EquipmentType
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

type CreateTennisRacketInput struct {
	Name    string   `json:"name" bson:"name"`
	Brand   *string  `json:"brand,omitempty" bson:"brand,omitempty"`
	BrandID *int     `json:"brandId,omitempty" bson:"brandId,omitempty"`
	Model   *string  `json:"model,omitempty" bson:"model,omitempty"`
	ModelID *int     `json:"modelId,omitempty" bson:"modelId,omitempty"`
	Weight  *float64 `json:"weight,omitempty" bson:"weight,omitempty"`
}

type CreateTennisStringInput struct {
	Name          string              `json:"name" bson:"name"`
	Brand         *string             `json:"brand,omitempty" bson:"brand,omitempty"`
	BrandID       *int                `json:"brandId,omitempty" bson:"brandId,omitempty"`
	Model         *string             `json:"model,omitempty" bson:"model,omitempty"`
	ModelID       *int                `json:"modelId,omitempty" bson:"modelId,omitempty"`
	Tension       *StringTensionInput `json:"tension,omitempty" bson:"tension,omitempty"`
	StringingDate *time.Time          `json:"stringingDate,omitempty" bson:"stringingDate,omitempty"`
	BurstDate     *time.Time          `json:"burstDate,omitempty" bson:"burstDate,omitempty"`
}

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

type StringTension struct {
	Mains   *int `json:"mains,omitempty" bson:"mains,omitempty"`
	Crosses *int `json:"crosses,omitempty" bson:"crosses,omitempty"`
}

type StringTensionInput struct {
	Mains   *int `json:"mains,omitempty" bson:"mains,omitempty"`
	Crosses *int `json:"crosses,omitempty" bson:"crosses,omitempty"`
}

type TennisRacket struct {
	ID              primitive.ObjectID  `json:"id" bson:"_id"`
	OwnerID         primitive.ObjectID  `json:"ownerId" bson:"ownerId"`
	Name            string              `json:"name" bson:"name"`
	Type            EquipmentType       `json:"type" bson:"type"`
	CreatedAt       time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt" bson:"updatedAt"`
	CurrentStringID *primitive.ObjectID `json:"currentStringId,omitempty" bson:"currentStringId,omitempty"`
	Brand           *string             `json:"brand,omitempty" bson:"brand,omitempty"`
	BrandID         *int                `json:"brandId,omitempty" bson:"brandId,omitempty"`
	Model           *string             `json:"model,omitempty" bson:"model,omitempty"`
	ModelID         *int                `json:"modelId,omitempty" bson:"modelId,omitempty"`
	Weight          *float64            `json:"weight,omitempty" bson:"weight,omitempty"`
}

func (TennisRacket) IsEquipment()                        {}
func (this TennisRacket) GetID() primitive.ObjectID      { return this.ID }
func (this TennisRacket) GetOwnerID() primitive.ObjectID { return this.OwnerID }
func (this TennisRacket) GetName() string                { return this.Name }
func (this TennisRacket) GetType() EquipmentType         { return this.Type }
func (this TennisRacket) GetCreatedAt() time.Time        { return this.CreatedAt }
func (this TennisRacket) GetUpdatedAt() time.Time        { return this.UpdatedAt }

func (TennisRacket) IsEntity() {}

type TennisString struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id"`
	OwnerID       primitive.ObjectID  `json:"ownerId" bson:"ownerId"`
	Name          string              `json:"name" bson:"name"`
	Type          EquipmentType       `json:"type" bson:"type"`
	CreatedAt     time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time           `json:"updatedAt" bson:"updatedAt"`
	Racket        *primitive.ObjectID `json:"racket,omitempty" bson:"racket,omitempty"`
	Brand         *string             `json:"brand,omitempty" bson:"brand,omitempty"`
	BrandID       *int                `json:"brandId,omitempty" bson:"brandId,omitempty"`
	Model         *string             `json:"model,omitempty" bson:"model,omitempty"`
	ModelID       *int                `json:"modelId,omitempty" bson:"modelId,omitempty"`
	Tension       *StringTension      `json:"tension,omitempty" bson:"tension,omitempty"`
	StringingDate *time.Time          `json:"stringingDate,omitempty" bson:"stringingDate,omitempty"`
	BurstDate     *time.Time          `json:"burstDate,omitempty" bson:"burstDate,omitempty"`
}

func (TennisString) IsEquipment()                        {}
func (this TennisString) GetID() primitive.ObjectID      { return this.ID }
func (this TennisString) GetOwnerID() primitive.ObjectID { return this.OwnerID }
func (this TennisString) GetName() string                { return this.Name }
func (this TennisString) GetType() EquipmentType         { return this.Type }
func (this TennisString) GetCreatedAt() time.Time        { return this.CreatedAt }
func (this TennisString) GetUpdatedAt() time.Time        { return this.UpdatedAt }

func (TennisString) IsEntity() {}

type UpdateTennisRacketInput struct {
	Name    *string  `json:"name,omitempty" bson:"name,omitempty"`
	Brand   *string  `json:"brand,omitempty" bson:"brand,omitempty"`
	BrandID *int     `json:"brandId,omitempty" bson:"brandId,omitempty"`
	Model   *string  `json:"model,omitempty" bson:"model,omitempty"`
	ModelID *int     `json:"modelId,omitempty" bson:"modelId,omitempty"`
	Weight  *float64 `json:"weight,omitempty" bson:"weight,omitempty"`
}

type UpdateTennisStringInput struct {
	Name          *string             `json:"name,omitempty" bson:"name,omitempty"`
	Brand         *string             `json:"brand,omitempty" bson:"brand,omitempty"`
	BrandID       *int                `json:"brandId,omitempty" bson:"brandId,omitempty"`
	Model         *string             `json:"model,omitempty" bson:"model,omitempty"`
	ModelID       *int                `json:"modelId,omitempty" bson:"modelId,omitempty"`
	Tension       *StringTensionInput `json:"tension,omitempty" bson:"tension,omitempty"`
	StringingDate *time.Time          `json:"stringingDate,omitempty" bson:"stringingDate,omitempty"`
	BurstDate     *time.Time          `json:"burstDate,omitempty" bson:"burstDate,omitempty"`
}

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	MyTennisRackets []*TennisRacket    `json:"myTennisRackets" bson:"myTennisRackets"`
	MyTennisStrings []*TennisString    `json:"myTennisStrings" bson:"myTennisStrings"`
}

func (User) IsEntity() {}

type EquipmentType string

const (
	EquipmentTypeTennisRacket EquipmentType = "TENNIS_RACKET"
	EquipmentTypeTennisString EquipmentType = "TENNIS_STRING"
)

var AllEquipmentType = []EquipmentType{
	EquipmentTypeTennisRacket,
	EquipmentTypeTennisString,
}

func (e EquipmentType) IsValid() bool {
	switch e {
	case EquipmentTypeTennisRacket, EquipmentTypeTennisString:
		return true
	}
	return false
}

func (e EquipmentType) String() string {
	return string(e)
}

func (e *EquipmentType) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EquipmentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EquipmentType", str)
	}
	return nil
}

func (e EquipmentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type StringGauge string

const (
	StringGaugeGauge15  StringGauge = "GAUGE_15"
	StringGaugeGauge15l StringGauge = "GAUGE_15L"
	StringGaugeGauge16  StringGauge = "GAUGE_16"
	StringGaugeGauge16l StringGauge = "GAUGE_16L"
	StringGaugeGauge17  StringGauge = "GAUGE_17"
	StringGaugeGauge18  StringGauge = "GAUGE_18"
	StringGaugeGauge19  StringGauge = "GAUGE_19"
)

var AllStringGauge = []StringGauge{
	StringGaugeGauge15,
	StringGaugeGauge15l,
	StringGaugeGauge16,
	StringGaugeGauge16l,
	StringGaugeGauge17,
	StringGaugeGauge18,
	StringGaugeGauge19,
}

func (e StringGauge) IsValid() bool {
	switch e {
	case StringGaugeGauge15, StringGaugeGauge15l, StringGaugeGauge16, StringGaugeGauge16l, StringGaugeGauge17, StringGaugeGauge18, StringGaugeGauge19:
		return true
	}
	return false
}

func (e StringGauge) String() string {
	return string(e)
}

func (e *StringGauge) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = StringGauge(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid StringGauge", str)
	}
	return nil
}

func (e StringGauge) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
