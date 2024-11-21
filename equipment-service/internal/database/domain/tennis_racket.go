package domain

import (
    "time"
	"github.com/CourtIQ/courtiq-backend/equipment-service//graph/model"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type TennisRacket struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Brand     string            `bson:"brand"`
    Model     string            `bson:"model"`
    HeadSize  float64           `bson:"head_size"`
    Weight    float64           `bson:"weight"`
    CreatedAt time.Time         `bson:"created_at"`
    UpdatedAt time.Time         `bson:"updated_at"`
}

func NewTennisRacket() *TennisRacket {
    return &TennisRacket{
        ID:        primitive.NewObjectID(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}