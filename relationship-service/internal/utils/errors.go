package utils

import "fmt"

// UIError is a custom error type that implements `error`.
type UIError struct {
	Code    string // e.g. "COACHSHIP_NOT_FOUND"
	Message string // e.g. "Coachship does not exist."
}

func (e UIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// You can define as many error variables as you need:
var (
	ErrCoachshipNotFound = UIError{
		Code:    "COACHSHIP_NOT_FOUND",
		Message: "Coachship does not exist.",
	}
	ErrMissingCoachshipID = UIError{
		Code:    "MISSING_COACHSHIP_ID",
		Message: "Cannot update a coachship without an ID.",
	}
	ErrMissingFriendshipID = UIError{
		Code:    "MISSING_FRIENDSHIP_ID",
		Message: "Cannot update a friendship without an ID.",
	}
	ErrFriendshipNotFound = UIError{
		Code:    "FRIENDSHIP_NOT_FOUND",
		Message: "Friendship does not exist.",
	}

	// ----------------------------------------------------------------------
	// Additional errors for forbidden access or permission-related issues.
	// ----------------------------------------------------------------------

	// ErrFriendshipForbidden is returned when a user tries to access or modify
	// a friendship they are not part of.
	ErrFriendshipForbidden = UIError{
		Code:    "FRIENDSHIP_FORBIDDEN",
		Message: "You do not have permission to view or modify this friendship.",
	}

	// ErrCoachshipForbidden is returned when a user tries to access or modify
	// a coachship they are not part of.
	ErrCoachshipForbidden = UIError{
		Code:    "COACHSHIP_FORBIDDEN",
		Message: "You do not have permission to view or modify this coachship.",
	}

	// ErrNotParticipantInRelationship can be used generically if you want a single
	// error for any relationship type, indicating the user is not a participant.
	ErrNotParticipantInRelationship = UIError{
		Code:    "RELATIONSHIP_NOT_PARTICIPANT",
		Message: "You are not a participant in this relationship.",
	}

	// ----------------------------------------------------------------------
	// Error for attempting to send a friend request to oneself.
	// ----------------------------------------------------------------------

	ErrCannotSendFriendRequestToSelf = UIError{
		Code:    "FRIENDSHIP_SELF_REQUEST",
		Message: "Cannot send a friend request to yourself.",
	}

	ErrCannotSendCoachRequestToSelf = UIError{
		Code:    "COACH_SELF_REQUEST",
		Message: "Cannot send a coach request to yourself.",
	}

	ErrCannotSendStudentRequestToSelf = UIError{
		Code:    "STUDENT_SELF_REQUEST",
		Message: "Cannot send a student request to yourself.",
	}

	// ----------------------------------------------------------------------
	// Error for indicating a friendship already exists between two users.
	// ----------------------------------------------------------------------

	ErrFriendshipAlreadyExists = UIError{
		Code:    "FRIENDSHIP_ALREADY_EXISTS",
		Message: "You are already a friend of this user or have a pending request.",
	}

	// ----------------------------------------------------------------------
	// Error indicating the user is already a coach or has a pending request.
	// ----------------------------------------------------------------------

	ErrAlreadyCoachOrPending = UIError{
		Code:    "COACHSHIP_ALREADY_EXISTS",
		Message: "You are already a coach for this user or have a pending request.",
	}

	// ----------------------------------------------------------------------
	// Error indicating the user is already a student or has a pending request.
	// ----------------------------------------------------------------------

	ErrAlreadyStudentOrPending = UIError{
		Code:    "COACHSHIP_ALREADY_EXISTS",
		Message: "You are already a student of this user or have a pending request.",
	}
)
