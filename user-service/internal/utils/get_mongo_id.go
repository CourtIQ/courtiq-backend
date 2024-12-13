package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

// UserClaims represents the JSON structure of the user claims
type UserClaims struct {
	MongoID       string `json:"mongoId"`
	Iss           string `json:"iss"`
	Aud           string `json:"aud"`
	AuthTime      int64  `json:"auth_time"`
	UserID        string `json:"user_id"`
	Sub           string `json:"sub"`
	Iat           int64  `json:"iat"`
	Exp           int64  `json:"exp"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Firebase      struct {
		Identities struct {
			Email []string `json:"email"`
		} `json:"identities"`
		SignInProvider string `json:"sign_in_provider"`
	} `json:"firebase"`
	UID string `json:"uid"`
}

// GetMongoID extracts the mongoId from the X-User-Claims header in the given request.
// Returns the mongoId if found, otherwise an error.
func GetMongoID(r *http.Request) (string, error) {
	claimsStr := r.Header.Get("X-User-Claims")
	if claimsStr == "" {
		return "", errors.New("no user claims found in X-User-Claims header")
	}

	var claims UserClaims
	if err := json.Unmarshal([]byte(claimsStr), &claims); err != nil {
		return "", err
	}

	if claims.MongoID == "" {
		return "", errors.New("mongoId not present in user claims")
	}

	return claims.MongoID, nil
}
