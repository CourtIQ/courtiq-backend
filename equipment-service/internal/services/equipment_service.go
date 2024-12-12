package services

import (
	"context"
	"fmt"
	"time"

	"github.com/CourtIQ/courtiq-backend/equipment-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/equipment-service/internal/repository"
	"github.com/CourtIQ/courtiq-backend/equipment-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EquipmentService struct {
	racketRepo repository.TennisRacketRepository
	stringRepo repository.TennisStringRepository
}

// Example AuthConfig. Adjust EnableAuth as needed.
var authCfg = utils.AuthConfig{EnableAuth: false}

// NewEquipmentService constructs an EquipmentService with the given repositories.
func NewEquipmentService(
	racketRepo repository.TennisRacketRepository,
	stringRepo repository.TennisStringRepository,
) EquipmentServiceIntf {
	return &EquipmentService{
		racketRepo: racketRepo,
		stringRepo: stringRepo,
	}
}

// applyPaginationDefaults sets default limit/offset values if nil
func applyPaginationDefaults(limit, offset *int) (int, int) {
	l, o := 10, 0
	if limit != nil {
		l = *limit
	}
	if offset != nil {
		o = *offset
	}
	return l, o
}

// ---- Mutations - Tennis Racket ----

func (s *EquipmentService) CreateTennisRacket(ctx context.Context, input model.CreateTennisRacketInput) (*model.TennisRacket, error) {
	uid, err := utils.GetUserIDFromContext(ctx, authCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Convert uid to ObjectID
	ownerID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	racket := &model.TennisRacket{
		ID:        primitive.NewObjectID(),
		OwnerID:   ownerID,
		Name:      input.Name,
		Brand:     input.Brand,
		Model:     input.Model,
		HeadSize:  input.HeadSize,
		Weight:    input.Weight,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Type:      model.EquipmentTypeTennisRacket,
	}

	if err := s.racketRepo.Insert(ctx, racket); err != nil {
		return nil, err
	}

	return racket, nil
}

func (s *EquipmentService) UpdateMyTennisRacket(ctx context.Context, id primitive.ObjectID, input model.UpdateTennisRacketInput) (*model.TennisRacket, error) {
	if _, err := utils.GetUserIDFromContext(ctx, authCfg); err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	racket, err := s.racketRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if racket == nil {
		return nil, nil // Not found
	}

	// Update fields directly since ownership check is removed
	if input.Name != nil {
		racket.Name = *input.Name
	}
	if input.Brand != nil {
		racket.Brand = input.Brand
	}
	if input.Model != nil {
		racket.Model = input.Model
	}
	if input.HeadSize != nil {
		racket.HeadSize = input.HeadSize
	}
	if input.Weight != nil {
		racket.Weight = input.Weight
	}
	racket.UpdatedAt = time.Now()

	if err := s.racketRepo.Update(ctx, racket); err != nil {
		return nil, err
	}
	return racket, nil
}

func (s *EquipmentService) DeleteMyTennisRacket(ctx context.Context, id primitive.ObjectID) (*model.TennisRacket, error) {
	if _, err := utils.GetUserIDFromContext(ctx, authCfg); err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	racket, err := s.racketRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if racket == nil {
		return nil, nil // Not found
	}

	if err := s.racketRepo.Delete(ctx, id); err != nil {
		return nil, err
	}

	return racket, nil
}

// ---- Mutations - Tennis String ----

func (s *EquipmentService) CreateTennisString(ctx context.Context, input model.CreateTennisStringInput) (*model.TennisString, error) {
	uid, err := utils.GetUserIDFromContext(ctx, authCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Convert uid to ObjectID
	ownerID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	str := &model.TennisString{
		ID:            primitive.NewObjectID(),
		OwnerID:       ownerID,
		Name:          input.Name,
		Brand:         input.Brand,
		Model:         input.Model,
		Gauge:         input.Gauge,
		Tension:       utils.ConvertStringTensionInput(input.Tension),
		StringingDate: input.StringingDate,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Type:          model.EquipmentTypeTennisString,
	}
	if err := s.stringRepo.Insert(ctx, str); err != nil {
		return nil, err
	}
	return str, nil
}

func (s *EquipmentService) UpdateMyTennisString(ctx context.Context, id primitive.ObjectID, input model.UpdateTennisStringInput) (*model.TennisString, error) {
	if _, err := utils.GetUserIDFromContext(ctx, authCfg); err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	strObj, err := s.stringRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if strObj == nil {
		return nil, nil
	}

	// Update fields directly (no ownership check)
	if input.Name != nil {
		strObj.Name = *input.Name
	}
	if input.Brand != nil {
		strObj.Brand = input.Brand
	}
	if input.Model != nil {
		strObj.Model = input.Model
	}
	if input.Gauge != nil {
		strObj.Gauge = input.Gauge
	}
	if input.Tension != nil {
		strObj.Tension = utils.ConvertStringTensionInput(input.Tension)
	}
	if input.StringingDate != nil {
		strObj.StringingDate = input.StringingDate
	}
	if input.BurstDate != nil {
		strObj.BurstDate = input.BurstDate
	}

	strObj.UpdatedAt = time.Now()

	if err := s.stringRepo.Update(ctx, strObj); err != nil {
		return nil, err
	}
	return strObj, nil
}

func (s *EquipmentService) DeleteMyTennisString(ctx context.Context, id primitive.ObjectID) (*model.TennisString, error) {
	if _, err := utils.GetUserIDFromContext(ctx, authCfg); err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	strObj, err := s.stringRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if strObj == nil {
		return nil, nil
	}

	if err := s.stringRepo.Delete(ctx, id); err != nil {
		return nil, err
	}
	return strObj, nil
}

// ---- Mutations - Racket-String Operations ----

func (s *EquipmentService) AssignStringToMyRacket(ctx context.Context, racketID primitive.ObjectID, stringID primitive.ObjectID) (*model.TennisRacket, error) {
	if _, err := utils.GetUserIDFromContext(ctx, authCfg); err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	racket, err := s.racketRepo.FindByID(ctx, racketID)
	if err != nil {
		return nil, err
	}
	if racket == nil {
		return nil, nil
	}

	strObj, err := s.stringRepo.FindByID(ctx, stringID)
	if err != nil {
		return nil, err
	}
	if strObj == nil {
		return nil, nil
	}

	racket.CurrentString = strObj
	racket.UpdatedAt = time.Now()
	if err := s.racketRepo.Update(ctx, racket); err != nil {
		return nil, err
	}
	return racket, nil
}

func (s *EquipmentService) RemoveStringFromMyRacket(ctx context.Context, racketID primitive.ObjectID) (*model.TennisRacket, error) {
	if _, err := utils.GetUserIDFromContext(ctx, authCfg); err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	racket, err := s.racketRepo.FindByID(ctx, racketID)
	if err != nil {
		return nil, err
	}
	if racket == nil {
		return nil, nil
	}

	racket.CurrentString = nil
	racket.UpdatedAt = time.Now()
	if err := s.racketRepo.Update(ctx, racket); err != nil {
		return nil, err
	}
	return racket, nil
}

// ---- Mutations - String Status Operations ----

func (s *EquipmentService) MarkMyStringAsBurst(ctx context.Context, stringID primitive.ObjectID, burstDate time.Time) (*model.TennisString, error) {
	if _, err := utils.GetUserIDFromContext(ctx, authCfg); err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	strObj, err := s.stringRepo.FindByID(ctx, stringID)
	if err != nil {
		return nil, err
	}
	if strObj == nil {
		return nil, nil
	}

	strObj.BurstDate = &burstDate
	strObj.UpdatedAt = time.Now()
	if err := s.stringRepo.Update(ctx, strObj); err != nil {
		return nil, err
	}
	return strObj, nil
}

func (s *EquipmentService) UpdateMyStringTension(ctx context.Context, stringID primitive.ObjectID, tension model.StringTensionInput) (*model.TennisString, error) {
	if _, err := utils.GetUserIDFromContext(ctx, authCfg); err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	strObj, err := s.stringRepo.FindByID(ctx, stringID)
	if err != nil {
		return nil, err
	}
	if strObj == nil {
		return nil, nil
	}

	strObj.Tension = &model.StringTension{
		Mains:   tension.Mains,
		Crosses: tension.Crosses,
	}
	strObj.UpdatedAt = time.Now()
	if err := s.stringRepo.Update(ctx, strObj); err != nil {
		return nil, err
	}
	return strObj, nil
}

// ---- Queries ----

func (s *EquipmentService) MyTennisRackets(ctx context.Context, limit *int, offset *int) ([]*model.TennisRacket, error) {
	uid, err := utils.GetUserIDFromContext(ctx, authCfg)
	if err != nil {
		return []*model.TennisRacket{}, fmt.Errorf("failed to get user ID: %w", err)
	}

	l, o := applyPaginationDefaults(limit, offset)
	allRackets, err := s.racketRepo.Find(ctx, bson.M{"ownerId": uid})
	if err != nil {
		return nil, err
	}

	end := o + l
	if end > len(allRackets) {
		end = len(allRackets)
	}
	if o > len(allRackets) {
		return []*model.TennisRacket{}, nil
	}

	return allRackets[o:end], nil
}

func (s *EquipmentService) MyTennisRacket(ctx context.Context, id primitive.ObjectID) (*model.TennisRacket, error) {
	if _, err := utils.GetUserIDFromContext(ctx, authCfg); err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	racket, err := s.racketRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if racket == nil {
		return nil, nil
	}
	return racket, nil
}

func (s *EquipmentService) MyTennisStrings(ctx context.Context, limit *int, offset *int) ([]*model.TennisString, error) {
	uid, err := utils.GetUserIDFromContext(ctx, authCfg)
	if err != nil {
		return []*model.TennisString{}, fmt.Errorf("failed to get user ID: %w", err)
	}

	l, o := applyPaginationDefaults(limit, offset)
	allStrings, err := s.stringRepo.Find(ctx, bson.M{"ownerId": uid})
	if err != nil {
		return nil, err
	}

	end := o + l
	if end > len(allStrings) {
		end = len(allStrings)
	}
	if o > len(allStrings) {
		return []*model.TennisString{}, nil
	}

	return allStrings[o:end], nil
}

func (s *EquipmentService) MyTennisString(ctx context.Context, id primitive.ObjectID) (*model.TennisString, error) {
	if _, err := utils.GetUserIDFromContext(ctx, authCfg); err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	strObj, err := s.stringRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if strObj == nil {
		return nil, nil
	}
	return strObj, nil
}

func (s *EquipmentService) MyEquipment(ctx context.Context, limit *int, offset *int) ([]model.Equipment, error) {
	uid, err := utils.GetUserIDFromContext(ctx, authCfg)
	if err != nil {
		return []model.Equipment{}, fmt.Errorf("failed to get user ID: %w", err)
	}

	l, o := applyPaginationDefaults(limit, offset)

	// Fetch rackets
	rackets, err := s.racketRepo.Find(ctx, bson.M{"ownerId": uid})
	if err != nil {
		return nil, err
	}

	// Fetch strings
	strings, err := s.stringRepo.Find(ctx, bson.M{"ownerId": uid})
	if err != nil {
		return nil, err
	}

	// Combine
	equipment := make([]model.Equipment, 0, len(rackets)+len(strings))
	for _, r := range rackets {
		equipment = append(equipment, r)
	}
	for _, s := range strings {
		equipment = append(equipment, s)
	}

	end := o + l
	if end > len(equipment) {
		end = len(equipment)
	}
	if o > len(equipment) {
		return []model.Equipment{}, nil
	}

	return equipment[o:end], nil
}

// ---- Entity Resolvers ----

func (s *EquipmentService) FindEquipmentByID(ctx context.Context, id primitive.ObjectID) (model.Equipment, error) {
	racket, err := s.racketRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if racket != nil {
		return racket, nil
	}

	strObj, err := s.stringRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if strObj != nil {
		return strObj, nil
	}

	return nil, nil // Not found
}

func (s *EquipmentService) FindTennisRacketByID(ctx context.Context, id primitive.ObjectID) (*model.TennisRacket, error) {
	return s.racketRepo.FindByID(ctx, id)
}

func (s *EquipmentService) FindTennisStringByID(ctx context.Context, id primitive.ObjectID) (*model.TennisString, error) {
	return s.stringRepo.FindByID(ctx, id)
}
