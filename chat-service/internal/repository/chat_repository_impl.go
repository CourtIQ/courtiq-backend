package repository

import (
	"context"
	"errors"

	"github.com/CourtIQ/courtiq-backend/chat-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
	sharedErrors "github.com/CourtIQ/courtiq-backend/shared/pkg/errors"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// chatRepository is the concrete implementation of ChatRepository interface
type chatRepository struct {
	repository.BaseRepository[model.Chat]
}

// NewChatRepository creates a new instance of ChatRepository using the provided factory
func NewChatRepository(factory *repository.RepositoryFactory) ChatRepository {
	baseRepo := repository.NewRepository[model.Chat](factory, db.ChatsCollection)
	return &chatRepository{
		BaseRepository: *baseRepo,
	}
}

// CreatePrivateChat creates a new private chat in the database
func (r *chatRepository) CreatePrivateChat(ctx context.Context, chat *model.PrivateChat) (*model.PrivateChat, error) {
	if chat == nil {
		return nil, sharedErrors.WrapError(errors.New("chat cannot be nil"), "invalid input")
	}
	if len(chat.ParticipantIds) != 2 {
		return nil, sharedErrors.WrapError(errors.New("private chat must have exactly two participants"), "invalid participant count")
	}
	if chat.ParticipantIds[0] == primitive.NilObjectID || chat.ParticipantIds[1] == primitive.NilObjectID {
		return nil, sharedErrors.WrapError(errors.New("participant IDs cannot be empty"), "invalid participant IDs")
	}
	if chat.ParticipantIds[0] == chat.ParticipantIds[1] {
		return nil, sharedErrors.WrapError(errors.New("participants cannot be the same"), "invalid participants")
	}
	if chat.Type != model.ChatTypePrivate {
		return nil, sharedErrors.WrapError(errors.New("invalid chat type"), "chat type must be private")
	}

	result, err := r.BaseRepository.Insert(ctx, chat)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, sharedErrors.WrapError(err, "private chat already exists")
		}
		return nil, sharedErrors.WrapError(err, "failed to create private chat")
	}

	createdChat, ok := result.(*model.PrivateChat)
	if !ok {
		return nil, sharedErrors.WrapError(errors.New("invalid type for inserted chat"), "failed to assert inserted chat type")
	}

	return createdChat, nil
}

// CreateGroupChat creates a new group chat in the database
func (r *chatRepository) CreateGroupChat(ctx context.Context, chat *model.GroupChat) (*model.GroupChat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")
}

// DeleteChat deletes a chat by its ID
func (r *chatRepository) DeleteChat(ctx context.Context, id *primitive.ObjectID) (*model.Chat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// GetChatByID retrieves a chat by its ID
func (r *chatRepository) GetChatByID(ctx context.Context, id *primitive.ObjectID) (*model.Chat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// GetPrivateChatByID retrieves a private chat by its ID
func (r *chatRepository) GetPrivateChatByID(ctx context.Context, id *primitive.ObjectID) (*model.PrivateChat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// GetGroupChatByID retrieves a group chat by its ID
func (r *chatRepository) GetGroupChatByID(ctx context.Context, id *primitive.ObjectID) (*model.GroupChat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// GetMyChats retrieves all chats for the current user with pagination
func (r *chatRepository) GetMyChats(ctx context.Context, limit int64, skip int64) ([]*model.Chat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// GetMyPrivateChats retrieves all private chats for the current user
func (r *chatRepository) GetMyPrivateChats(ctx context.Context) ([]*model.PrivateChat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// GetMyGroupChats retrieves all group chats for the current user
func (r *chatRepository) GetMyGroupChats(ctx context.Context) ([]*model.GroupChat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// UpdateGroupChat updates an existing group chat
func (r *chatRepository) UpdateGroupChat(ctx context.Context, chat *model.GroupChat) (*model.GroupChat, error) {
	return nil, sharedErrors.WrapError(errors.New("not implemented"), "AddImageMessage is not yet implemented")

}

// DoesPrivateChatExist checks if a private chat exists between two users
func (r *chatRepository) DoesPrivateChatExist(ctx context.Context, userID1, userID2 primitive.ObjectID) (bool, error) {
	return false, nil

}
