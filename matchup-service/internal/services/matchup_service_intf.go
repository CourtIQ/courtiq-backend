package services

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MatchUpServiceIntf defines the interface for matchup service operations
type MatchUpServiceIntf interface {
	// MatchUp format operations
	GetMatchUpFormats(ctx context.Context, limit *int, offset *int) ([]*model.MatchUpFormat, error)
	CreateMatchUpFormat(ctx context.Context, input model.MatchUpFormatInput) (*model.MatchUpFormat, error)
	
	// MatchUp operations
	InitiateMatchUp(ctx context.Context, input model.InitiateMatchUpInput) (*model.MatchUp, error)
	GetMatchUps(ctx context.Context, limit *int, offset *int) ([]*model.MatchUp, error)
	GetMatchUpById(ctx context.Context, id primitive.ObjectID) (*model.MatchUp, error)
	UpdateMatchUp(ctx context.Context, id primitive.ObjectID, status model.MatchUpStatus) (*model.MatchUp, error)
	
	// MatchUp shot operations
	AddShot(ctx context.Context, input model.AddShotInput) (*model.MatchUpShot, error)
	GetMatchUpShots(ctx context.Context, matchUpId primitive.ObjectID, limit *int, offset *int) ([]*model.MatchUpShot, error)
	GetShotsByGame(ctx context.Context, matchUpId primitive.ObjectID, setNumber int, gameNumber int) ([]*model.MatchUpShot, error)
}