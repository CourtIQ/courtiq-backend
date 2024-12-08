// satisfies.go (or the file where SatisfiesDirective is defined)
package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/CourtIQ/courtiq-backend/relationship-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/relationship-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// These should be set up in main or resolver initialization:
var RelationshipRepo repository.RelationshipRepository
var GetCurrentUserID func(ctx context.Context) (string, error)

func SatisfiesDirective(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
	relationshipStatus *model.RelationshipStatus,
	relationshipType *model.RelationshipType,
	requireParticipant *bool,
	requireSender *bool,
	requireReceiver *bool,
	noExistingFriendship *bool,
) (interface{}, error) {

	participant := requireParticipant != nil && *requireParticipant
	sender := requireSender != nil && *requireSender
	receiver := requireReceiver != nil && *requireReceiver
	noExisting := noExistingFriendship != nil && *noExistingFriendship

	// Get the current user ID
	currentUserID, err := GetCurrentUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	fc := graphql.GetFieldContext(ctx)

	// Extract arguments (if any)
	friendshipID, _ := fc.Args["friendshipId"].(string)
	receiverID, _ := fc.Args["receiverId"].(string)

	// Handle noExistingFriendship scenario first
	if noExisting {
		if receiverID == "" {
			return nil, fmt.Errorf("noExistingFriendship requires receiverId argument")
		}

		// Use map[string]interface{} consistently
		noExistFilter := map[string]interface{}{
			"type": "FRIENDSHIP",
			"participantIds": map[string]interface{}{
				"$all": []string{currentUserID, receiverID},
			},
			"status": map[string]interface{}{
				"$in": []string{"PENDING", "ACTIVE"},
			},
		}

		count, err := RelationshipRepo.Count(ctx, noExistFilter)
		if err != nil {
			return nil, fmt.Errorf("failed checking noExistingFriendship: %w", err)
		}
		if count > 0 {
			return nil, fmt.Errorf("forbidden: existing friendship or pending request already present")
		}

		// If no other conditions specified, just proceed
		if relationshipStatus == nil && relationshipType == nil && !participant && !sender && !receiver {
			return next(ctx)
		}
	}

	// If we need to check other conditions (status/type/participant/sender/receiver),
	// we must have a known relationship doc identified by friendshipId
	if relationshipStatus != nil || relationshipType != nil || participant || sender || receiver {
		if friendshipID == "" {
			return nil, fmt.Errorf("missing friendshipId for relationship checks")
		}

		objID, err := primitive.ObjectIDFromHex(friendshipID)
		if err != nil {
			return nil, fmt.Errorf("invalid friendshipId: %w", err)
		}

		filter := map[string]interface{}{
			"_id": objID,
		}

		// Apply status condition if set
		if relationshipStatus != nil {
			filter["status"] = relationshipStatus.String()
		}
		// Apply type condition if set
		if relationshipType != nil {
			filter["type"] = relationshipType.String()
		}
		// If requireParticipant, ensure currentUser in participantIds
		if participant {
			filter["participantIds"] = currentUserID
		}
		// If requireSender, ensure currentUser is sender
		if sender {
			filter["requesterId"] = currentUserID
		}
		// If requireReceiver, ensure currentUser is receiver
		if receiver {
			filter["receiverId"] = currentUserID
		}

		count, err := RelationshipRepo.Count(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to check conditions: %w", err)
		}

		if count == 0 {
			return nil, fmt.Errorf("forbidden: conditions not met")
		}
	}

	// Conditions satisfied
	return next(ctx)
}
