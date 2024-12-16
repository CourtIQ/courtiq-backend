package services

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/equipment-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EquipmentServiceIntf defines a 1:1 mapping of all resolver methods to service methods
type EquipmentServiceIntf interface {
	// Mutations - Tennis Racket
	CreateTennisRacket(ctx context.Context, input model.CreateTennisRacketInput) (*model.TennisRacket, error)
	UpdateMyTennisRacket(ctx context.Context, id primitive.ObjectID, input model.UpdateTennisRacketInput) (*model.TennisRacket, error)
	DeleteMyTennisRacket(ctx context.Context, id primitive.ObjectID) (bool, error)

	// Mutations - Tennis String
	CreateTennisString(ctx context.Context, input model.CreateTennisStringInput) (*model.TennisString, error)
	UpdateMyTennisString(ctx context.Context, id primitive.ObjectID, input model.UpdateTennisStringInput) (*model.TennisString, error)
	DeleteMyTennisString(ctx context.Context, id primitive.ObjectID) (bool, error)

	// Mutations - Racket-String Operations
	AssignStringToMyRacket(ctx context.Context, racketID primitive.ObjectID, stringID primitive.ObjectID) (*model.TennisRacket, error)

	// Queries (now with pagination parameters for list queries)
	MyTennisRackets(ctx context.Context, limit *int, offset *int) ([]*model.TennisRacket, error)
	MyTennisRacket(ctx context.Context, id primitive.ObjectID) (*model.TennisRacket, error)
	MyTennisStrings(ctx context.Context, limit *int, offset *int) ([]*model.TennisString, error)
	MyTennisString(ctx context.Context, id primitive.ObjectID) (*model.TennisString, error)
	MyEquipment(ctx context.Context, limit *int, offset *int) ([]model.Equipment, error)

	// Entity Resolvers
	FindEquipmentByID(ctx context.Context, id primitive.ObjectID) (model.Equipment, error)
	FindTennisRacketByID(ctx context.Context, id primitive.ObjectID) (*model.TennisRacket, error)
	FindTennisStringByID(ctx context.Context, id primitive.ObjectID) (*model.TennisString, error)
}
