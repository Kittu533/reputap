package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/kittu533/reputap-be/internal/response"
)

// 💡 Go Concept: Context Keys
// To store values in context, we use a custom type as key.
// This prevents collisions with other packages that might use
// the same string key. It's a Go best practice.
type contextKey string

const claimsKey contextKey = "claims"

// Middleware returns an HTTP middleware that validates JWT tokens.
//
// 💡 Go Concept: Middleware pattern
// A middleware is a function that wraps an http.Handler.
// It runs BEFORE the actual handler, and can:
// - Reject the request (return early without calling next)
// - Modify the request (add values to context)
// - Modify the response (add headers)
//
// The pattern is: func(next http.Handler) http.Handler
// Chi router uses this pattern natively with r.Use().
func Middleware(service *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Step 1: Extract token from Authorization header
			// Format: "Bearer <token>"
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "Authorization header is required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.Unauthorized(w, "Invalid authorization format. Use: Bearer <token>")
				return
			}

			tokenString := parts[1]

			// Step 2: Validate token and extract claims
			claims, err := service.ValidateToken(tokenString)
			if err != nil {
				response.Unauthorized(w, "Invalid or expired token")
				return
			}

			// Step 3: Store claims in request context
			// 💡 Go Concept: context.WithValue
			// This creates a new context with the claims embedded.
			// Downstream handlers can retrieve it with GetClaimsFromContext().
			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequirePermission returns middleware that checks if the user has a specific permission.
//
// 💡 Go Concept: Closure
// This function returns another function (middleware).
// The returned function "closes over" the permission and service variables,
// meaning it can access them even though they're defined in the outer function.
func RequirePermission(service *Service, permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaimsFromContext(r.Context())
			if !ok {
				response.Unauthorized(w, "Not authenticated")
				return
			}

			// Super admin bypasses permission checks
			if claims.RoleName == "super_admin" {
				next.ServeHTTP(w, r)
				return
			}

			permissions, err := service.GetUserPermissions(r.Context(), claims.RoleID)
			if err != nil {
				response.InternalError(w)
				return
			}

			for _, p := range permissions {
				if p == permission {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.Forbidden(w, "You don't have permission to perform this action")
		})
	}
}

// GetClaimsFromContext retrieves the JWT claims from the request context.
//
// 💡 Go Concept: Multiple return values
// Go functions can return multiple values. The second bool value
// follows the "comma ok" idiom — it tells you if the value exists.
// This is similar to map lookups: val, ok := myMap[key]
func GetClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	return claims, ok
}
