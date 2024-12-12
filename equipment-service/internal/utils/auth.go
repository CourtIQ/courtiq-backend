package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
)

// UserContextKey is a custom type used as the key for storing user data in context.
type userContextKey struct{}

type AuthConfig struct {
	EnableAuth bool
}

// NewUserMiddleware creates a middleware that extracts user information from the "x-user" header.
// If auth is disabled, GetUserIDFromContext will return a test value instead.
func NewUserMiddleware(cfg AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if cfg.EnableAuth {
				xUser := r.Header.Get("x-user")
				if xUser != "" {
					decoded, err := base64.StdEncoding.DecodeString(xUser)
					if err == nil {
						var data map[string]interface{}
						if err := json.Unmarshal(decoded, &data); err == nil {
							ctx = context.WithValue(ctx, userContextKey{}, data)
						}
					}
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext retrieves the user ID ("uid") from the context.
// If auth is disabled, it returns a dummy UID.
// If auth is enabled and no user data/uid is found, it returns an error.
func GetUserIDFromContext(ctx context.Context, cfg AuthConfig) (string, error) {
	if !cfg.EnableAuth {
		// Auth disabled, return dummy UID without error
		return "67575a5b47d020255890ea63", nil
	}

	data, ok := ctx.Value(userContextKey{}).(map[string]interface{})
	if !ok {
		return "", errors.New("no user data in context")
	}

	uidVal, exists := data["uid"]
	if !exists {
		return "", errors.New("no uid in user data")
	}

	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		return "", errors.New("uid is not a string or is empty")
	}

	return uid, nil
}
