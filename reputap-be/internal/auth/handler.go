package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/kittu533/reputap-be/internal/response"
)

// Handler handles HTTP requests for authentication.
//
// 💡 Go Concept: HTTP Handlers
// In Go, an HTTP handler is any function with signature:
//   func(w http.ResponseWriter, r *http.Request)
//
// - w (ResponseWriter): you write the HTTP response to this
// - r (*Request): contains the incoming HTTP request data
//
// We group related handlers in a struct so they can share dependencies.
type Handler struct {
	service  *Service
	validate *validator.Validate
}

// NewHandler creates a new auth handler.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service:  service,
		validate: validator.New(),
	}
}

// Login handles POST /api/v1/auth/login
//
// 💡 Go Concept: JSON Decoding
// json.NewDecoder(r.Body).Decode(&req) reads the request body
// and unmarshals the JSON into the struct. The struct tags like
// `json:"email"` tell the decoder which JSON key maps to which field.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Decode JSON body into struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	// Validate struct fields using validator tags
	// 💡 Go Concept: Struct validation
	// The validator package reads struct tags like `validate:"required,email"`
	// and checks if the values conform. This is similar to Zod in TypeScript.
	if err := h.validate.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var details []response.FieldError
			for _, e := range validationErrors {
				details = append(details, response.FieldError{
					Field:   e.Field(),
					Message: formatValidationError(e),
				})
			}
			response.ValidationErr(w, details)
			return
		}
		response.Err(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
		return
	}

	// Call service layer
	result, err := h.service.Login(r.Context(), req)
	if err != nil {
		// 💡 Go Concept: errors.Is()
		// We check which specific error was returned and map it
		// to the appropriate HTTP status code. This is how Go
		// separates business logic errors from HTTP concerns.
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			response.Unauthorized(w, "Invalid email or password")
		case errors.Is(err, ErrUserInactive):
			response.Forbidden(w, "Your account is inactive")
		default:
			response.InternalError(w)
		}
		return
	}

	response.JSON(w, http.StatusOK, result)
}

// Me handles GET /api/v1/auth/me
//
// 💡 Go Concept: Context values
// The auth middleware stores the user claims in the request context.
// We retrieve them here using r.Context().Value().
// This is how data flows from middleware to handlers in Go.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "Not authenticated")
		return
	}

	user, err := h.service.GetCurrentUser(r.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			response.Unauthorized(w, "User not found")
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// Refresh handles POST /api/v1/auth/refresh
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.Err(w, http.StatusBadRequest, "VALIDATION_ERROR", "Refresh token is required")
		return
	}

	result, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidToken):
			response.Unauthorized(w, "Invalid or expired refresh token")
		case errors.Is(err, ErrUserInactive):
			response.Forbidden(w, "Your account is inactive")
		default:
			response.InternalError(w)
		}
		return
	}

	response.JSON(w, http.StatusOK, result)
}

// Logout handles POST /api/v1/auth/logout
//
// 💡 Note: With JWT, logout is client-side (discard the token).
// For true server-side logout, you'd need a token blacklist — overkill for MVP.
// The client simply deletes the stored tokens.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

// formatValidationError converts a validator.FieldError to a human-readable message.
func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "email":
		return e.Field() + " must be a valid email"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters"
	case "max":
		return e.Field() + " must be at most " + e.Param() + " characters"
	default:
		return e.Field() + " is invalid"
	}
}
