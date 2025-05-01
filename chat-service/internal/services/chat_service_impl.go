package services

import (
	"context"
	"errors"
	"time"

	"github.com/CourtIQ/courtiq-backend/chat-service/graph/model"
	sharedErrors "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreatePrivateChat creates a new private chat between two users.
func (s *ChatServiceImpl) CreatePrivateChat(ctx context.Context, id primitive.ObjectID) (*model.PrivateChat, error) {
	currentUserID, err := middleware.GetMongoIDFromContext(ctx)
	if err != nil {
		return nil, sharedErrors.WrapError(err, "failed to get user ID from context")
	}

	// Validate participant ID
	if id == primitive.NilObjectID {
		return nil, sharedErrors.WrapError(errors.New("invalid participant ID"), "participant ID cannot be empty")
	}
	if id == currentUserID {
		return nil, sharedErrors.WrapError(errors.New("invalid participant"), "cannot create private chat with self")
	}

	newChat := &model.PrivateChat{
		ID:             primitive.NewObjectID(),
		ParticipantIds: []primitive.ObjectID{currentUserID, id},
		Type:           model.ChatTypePrivate,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	addedChat, err := s.chatRepo.CreatePrivateChat(ctx, newChat)

	if err != nil {
		return nil, sharedErrors.WrapError(err, "failed to create private chat")
	}

	return addedChat, nil
}

// CreateGroupChat creates a new group chat with the specified name and participants.
func (s *ChatServiceImpl) CreateGroupChat(ctx context.Context, name string, participantIDs []primitive.ObjectID) (*model.GroupChat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// AddParticipantsToGroupChat adds participants to an existing group chat.
func (s *ChatServiceImpl) AddParticipantsToGroupChat(ctx context.Context, chatID primitive.ObjectID, participantIDs []primitive.ObjectID) (*model.GroupChat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// RemoveParticipantsFromGroupChat removes participants from a group chat.
func (s *ChatServiceImpl) RemoveParticipantsFromGroupChat(ctx context.Context, chatID primitive.ObjectID, participantIDs []primitive.ObjectID) (*model.GroupChat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// LeaveGroupChat allows a user to leave a group chat.
func (s *ChatServiceImpl) LeaveGroupChat(ctx context.Context, chatID primitive.ObjectID) (*model.GroupChat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// GetChatByID retrieves a chat by its ID.
func (s *ChatServiceImpl) GetChatByID(ctx context.Context, chatID primitive.ObjectID) (model.Chat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// GetUserChats retrieves all chats for a user with pagination.
func (s *ChatServiceImpl) GetUserChats(ctx context.Context, limit *int, skip *int) ([]model.Chat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// SendTextMessage sends a text message to a chat.
func (s *ChatServiceImpl) SendTextMessage(ctx context.Context, chatID primitive.ObjectID, text string) (*model.TextMessage, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// GetMessagesByChatID retrieves messages for a chat with pagination.
func (s *ChatServiceImpl) GetMessagesByChatID(ctx context.Context, chatID primitive.ObjectID, limit int, skip int) ([]model.Message, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")
}
