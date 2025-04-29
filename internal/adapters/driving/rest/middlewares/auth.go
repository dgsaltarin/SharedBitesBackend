package middleware

import (
	"context"
	"log" // Use your preferred logger
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

// contextKey is a type used for context keys to avoid collisions.
type contextKey string

// UserContextKey is the key used to store authenticated user info in the context.
const UserContextKey contextKey = "authenticatedUser"

// AuthenticatedUser holds information about the verified user.
// Add more fields (Email, Name) if you extract them and need them downstream.
type AuthenticatedUser struct {
	UID string // Firebase User ID
	// Email string
	// Name string
}

// FirebaseAuthMiddleware creates a middleware handler that verifies Firebase ID tokens.
func FirebaseAuthMiddleware(authClient *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("Auth Middleware: Missing Authorization header") // Use logger
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			// Expecting "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Println("Auth Middleware: Invalid Authorization header format") // Use logger
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			idToken := parts[1]

			// Verify the ID token
			token, err := authClient.VerifyIDToken(r.Context(), idToken)
			if err != nil {
				// Log the specific error for debugging
				log.Printf("Auth Middleware: Error verifying Firebase ID token: %v", err)
				// Return a generic error to the client
				http.Error(w, "Invalid or expired authentication token", http.StatusUnauthorized)
				return
			}

			// Token is valid, extract user info and add to context
			authUser := AuthenticatedUser{
				UID: token.UID,
				// You can access other claims like:
				// Email: token.Claims["email"].(string), // Requires type assertion and error checking
				// Name: token.Claims["name"].(string),
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserContextKey, authUser)

			// Call the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves the authenticated user from the request context.
// Returns the user and a boolean indicating if the user was found.
func GetUserFromContext(ctx context.Context) (AuthenticatedUser, bool) {
	user, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	return user, ok
}

// RequireAuth is a helper function for handlers to easily get the user or write an error.
// DEPRECATED: It's generally better to let the middleware handle the 401/403 for missing/invalid auth,
// and have handlers assume the user exists if the middleware passed. This function remains
// primarily for checking if the context value was set correctly (defensive programming).
func RequireAuth(w http.ResponseWriter, r *http.Request) (AuthenticatedUser, bool) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		// This *shouldn't* happen if the middleware is applied correctly
		log.Println("Error: Authenticated user not found in context where expected")                    // Use logger
		http.Error(w, "Authentication required (user context missing)", http.StatusInternalServerError) // Or 401/403?
		return AuthenticatedUser{}, false
	}
	return user, true
}
