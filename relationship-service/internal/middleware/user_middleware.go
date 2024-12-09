// internal/middleware/user_middleware.go
package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
)

// Toggle this boolean to enable/disable auth header parsing.
// If disabled, GetUserIDFromContext always returns "testtter".
var enableAuth = false // Set to true to enable actual auth parsing

type userContextKey struct{}

// UserMiddleware extracts user information from the "x-user" header if auth is enabled.
// If disabled, it does nothing and user data will not be set, causing GetUserIDFromContext to return "testtter".
func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if enableAuth {
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

// GetUserIDFromContext retrieves the user ID ("uid") from the context.
// If auth is disabled, it always returns "testtter".
// If auth is enabled and no user data/uid is found, it returns an error.
func GetUserIDFromContext(ctx context.Context) (string, error) {
	if !enableAuth {
		// Auth disabled, return dummy UID without error
		return "sdddfdfsd", nil
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
