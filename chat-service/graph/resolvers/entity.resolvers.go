package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.61

import (
	"context"
	"fmt"

	"github.com/CourtIQ/courtiq-backend/chat-service/graph"
	"github.com/CourtIQ/courtiq-backend/chat-service/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindChatByID is the resolver for the findChatByID field.
func (r *entityResolver) FindChatByID(ctx context.Context, id primitive.ObjectID) (model.Chat, error) {
	panic(fmt.Errorf("not implemented: FindChatByID - findChatByID"))
}

// FindGroupChatByID is the resolver for the findGroupChatByID field.
func (r *entityResolver) FindGroupChatByID(ctx context.Context, id primitive.ObjectID) (*model.GroupChat, error) {
	panic(fmt.Errorf("not implemented: FindGroupChatByID - findGroupChatByID"))
}

// FindImageMessageByID is the resolver for the findImageMessageByID field.
func (r *entityResolver) FindImageMessageByID(ctx context.Context, id primitive.ObjectID) (*model.ImageMessage, error) {
	panic(fmt.Errorf("not implemented: FindImageMessageByID - findImageMessageByID"))
}

// FindMessageByID is the resolver for the findMessageByID field.
func (r *entityResolver) FindMessageByID(ctx context.Context, id primitive.ObjectID) (model.Message, error) {
	panic(fmt.Errorf("not implemented: FindMessageByID - findMessageByID"))
}

// FindPrivateChatByID is the resolver for the findPrivateChatByID field.
func (r *entityResolver) FindPrivateChatByID(ctx context.Context, id primitive.ObjectID) (*model.PrivateChat, error) {
	panic(fmt.Errorf("not implemented: FindPrivateChatByID - findPrivateChatByID"))
}

// FindTextMessageByID is the resolver for the findTextMessageByID field.
func (r *entityResolver) FindTextMessageByID(ctx context.Context, id primitive.ObjectID) (*model.TextMessage, error) {
	panic(fmt.Errorf("not implemented: FindTextMessageByID - findTextMessageByID"))
}

// Entity returns graph.EntityResolver implementation.
func (r *Resolver) Entity() graph.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
