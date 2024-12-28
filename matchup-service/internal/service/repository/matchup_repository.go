package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/matchup-service/graph/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type MatchUpRepository interface {
	CreateMatchUp(ctx context.Context, matchup *model.MatchUp) (*model.MatchUp, error)
	FindByID(ctx context.Context, id string) (*model.MatchUp, error)
}

type matchupRepository struct {
	pointsColl  *mongo.Collection
	matchupColl *mongo.Collection
}

func NewMatchUpRepository(mdb *db.MongoDB) MatchUpRepository {
	return &matchupRepository{
		pointsColl:  mdb.Points
		matchupColl: mdb.GetCollection("matchups"),
	}
}
