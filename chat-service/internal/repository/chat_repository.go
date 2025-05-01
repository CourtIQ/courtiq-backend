package repository

import (
	"context"

	"github.com/CourtIQ/courtiq-backend/chat-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatRepository defines the interface for chat-related database operations
type ChatRepository interface {
	repository.Repository[model.Chat]
	CreatePrivateChat(ctx context.Context, chat *model.PrivateChat) (*model.PrivateChat, error)
	CreateGroupChat(ctx context.Context, chat *model.GroupChat) (*model.GroupChat, error)

	DeleteChat(ctx context.Context, id *primitive.ObjectID) (*model.Chat, error)

	GetChatByID(ctx context.Context, id *primitive.ObjectID) (*model.Chat, error)
	GetPrivateChatByID(ctx context.Context, id *primitive.ObjectID) (*model.PrivateChat, error)
	GetGroupChatByID(ctx context.Context, id *primitive.ObjectID) (*model.GroupChat, error)

	GetMyChats(ctx context.Context, limit int64, skip int64) ([]*model.Chat, error)
	GetMyPrivateChats(ctx context.Context) ([]*model.PrivateChat, error)
	GetMyGroupChats(ctx context.Context) ([]*model.GroupChat, error)

	UpdateGroupChat(ctx context.Context, chat *model.GroupChat) (*model.GroupChat, error)

	DoesPrivateChatExist(ctx context.Context, userID1, userID2 primitive.ObjectID) (bool, error)
}
