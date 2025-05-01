package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/chat-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MessageRepository defines the interface for message-related database operations
type MessageRepository interface {
	AddTextMessage(ctx context.Context, message *model.TextMessage) (*model.TextMessage, error)
	AddImageMessage(ctx context.Context, message *model.ImageMessage) (*model.ImageMessage, error)
	DeleteMessage(ctx context.Context, messageID primitive.ObjectID) (*model.Message, error)
	DeleteImageMessage(ctx context.Context, message *model.ImageMessage) (*model.ImageMessage, error)
	UpdateTextMessage(ctx context.Context, chatID primitive.ObjectID, message *model.TextMessage) (*model.TextMessage, error)
	GetMessagesByChatID(ctx context.Context, chatID primitive.ObjectID, limit int64, skip int64) ([]*model.Message, error)
}
