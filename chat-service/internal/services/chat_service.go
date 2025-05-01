package services

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/chat-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/chat-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatService defines the interface for chat and message related operations.
type ChatService interface {
	// Chat Mutations
	CreatePrivateChat(ctx context.Context, id primitive.ObjectID) (*model.PrivateChat, error)
	CreateGroupChat(ctx context.Context, name string, participantIDs []primitive.ObjectID) (*model.GroupChat, error)
	AddParticipantsToGroupChat(ctx context.Context, chatID primitive.ObjectID, participantIDs []primitive.ObjectID) (*model.GroupChat, error)
	RemoveParticipantsFromGroupChat(ctx context.Context, chatID primitive.ObjectID, participantIDs []primitive.ObjectID) (*model.GroupChat, error)
	LeaveGroupChat(ctx context.Context, chatID primitive.ObjectID) (*model.GroupChat, error)

	// Chat Queries
	GetChatByID(ctx context.Context, chatID primitive.ObjectID) (model.Chat, error)
	GetUserChats(ctx context.Context, limit *int, skip *int) ([]model.Chat, error)

	// Message Mutations
	SendTextMessage(ctx context.Context, chatID primitive.ObjectID, text string) (*model.TextMessage, error)

	// Message Queries
	GetMessagesByChatID(ctx context.Context, chatID primitive.ObjectID, limit int, skip int) ([]model.Message, error)
}

// RelationshipService implements RelationshipServiceIntf
type ChatServiceImpl struct {
	chatRepo    repository.ChatRepository
	messageRepo repository.MessageRepository
}

// NewRelationshipService creates a new RelationshipService
func NewChatService(chatRepo repository.ChatRepository, messageRepo repository.MessageRepository) *ChatServiceImpl {
	return &ChatServiceImpl{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
	}
}
