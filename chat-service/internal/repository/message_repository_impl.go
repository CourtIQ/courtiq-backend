package repository

import (
	"context"
	"errors"

	"github.com/CourtIQ/courtiq-backend/chat-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	sharedErrors "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// messageRepository is the concrete implementation of MessageRepository interface
type messageRepository struct {
	repository.BaseRepository[model.Message]
}

// NewMessageRepository creates a new instance of MessageRepository using the provided factory
func NewMessageRepository(factory *repository.RepositoryFactory) MessageRepository {
	baseRepo := repository.NewRepository[model.Message](factory, db.MessagesCollection)
	return &messageRepository{
		BaseRepository: *baseRepo,
	}
}

// AddTextMessage adds a new text message to the database
func (r *messageRepository) AddTextMessage(ctx context.Context, message *model.TextMessage) (*model.TextMessage, error) {
	// Validate the incoming message
	if message == nil {
		return nil, sharedErrors.WrapError(errors.New("message is nil"), "invalid input")
	}
	if message.ChatID.IsZero() {
		return nil, sharedErrors.WrapError(errors.New("chat ID is not set"), "invalid input")
	}
	if message.SenderID.IsZero() {
		return nil, sharedErrors.WrapError(errors.New("sender ID is not set"), "invalid input")
	}
	if message.Text == "" {
		return nil, sharedErrors.WrapError(errors.New("text content is empty"), "invalid input")
	}
	if message.Type != model.MessageTypeText {
		return nil, sharedErrors.WrapError(errors.New("message type is not TEXT"), "invalid input")
	}

	// Verify that the sender ID matches the user ID from context for access control
	senderIDFromCtx, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, sharedErrors.WrapError(err, "failed to get sender ID from context")
	}
	if senderIDFromCtx != message.SenderID {
		return nil, sharedErrors.WrapError(errors.New("sender ID does not match context user ID"), "access denied")
	}

	// Use the embedded BaseRepository's Insert method
	insertedEntity, err := r.BaseRepository.Insert(ctx, message)
	if err != nil {
		return nil, sharedErrors.WrapError(err, "failed to insert text message")
	}

	// Type assert the inserted entity back to *model.TextMessage
	insertedMessage, ok := insertedEntity.(*model.TextMessage)
	if !ok {
		return nil, sharedErrors.WrapError(errors.New("inserted entity is not a TextMessage"), "unexpected type returned from insert")
	}

	return insertedMessage, nil
}

// AddImageMessage adds a new image message to the database (placeholder for future implementation)
func (r *messageRepository) AddImageMessage(ctx context.Context, message *model.ImageMessage) (*model.ImageMessage, error) {
	// TODO: Implement similar validation and insertion logic as AddTextMessage
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")
}

// DeleteMessage deletes a message by its ID (placeholder for future implementation)
func (r *messageRepository) DeleteMessage(ctx context.Context, messageID primitive.ObjectID) (*model.Message, error) {
	// TODO: Implement deletion with access control checks
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "DeleteMessage is not yet implemented")
}

// DeleteImageMessage deletes an image message (placeholder for future implementation)
func (r *messageRepository) DeleteImageMessage(ctx context.Context, message *model.ImageMessage) (*model.ImageMessage, error) {
	// TODO: Implement deletion with access control checks
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "DeleteImageMessage is not yet implemented")
}

// UpdateTextMessage updates a text message (placeholder for future implementation)
func (r *messageRepository) UpdateTextMessage(ctx context.Context, chatID primitive.ObjectID, message *model.TextMessage) (*model.TextMessage, error) {
	// TODO: Implement update with access control checks
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "UpdateTextMessage is not yet implemented")
}

// GetMessagesByChatID retrieves messages for a specific chat with pagination
func (r *messageRepository) GetMessagesByChatID(ctx context.Context, chatID primitive.ObjectID, limit int64, skip int64) ([]*model.Message, error) {
	// Validate input
	if chatID.IsZero() {
		return nil, sharedErrors.WrapError(errors.New("chat ID is not set"), "invalid input")
	}

	// Define filter for messages belonging to the chat
	filter := primitive.M{
		"chatId": chatID,
	}

	// Set pagination options
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(primitive.M{"createdAt": -1}) // Sort by creation time, descending (newest first)

	// Use BaseRepository's Find method to retrieve messages
	messages, err := r.BaseRepository.Find(ctx, filter, opts)
	if err != nil {
		return nil, sharedErrors.WrapError(err, "failed to retrieve messages for chat")
	}

	return messages, nil
}
