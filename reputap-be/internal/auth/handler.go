package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/kittu533/reputap-be/internal/response"
)

// Handler handles HTTP requests for authentication.
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
//	@Summary		Login
//	@Description	Authenticate admin user with email and password, returns JWT tokens
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginRequest	true	"Login credentials"
//	@Success		200		{object}	response.Response{data=TokenResponse}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		401		{object}	response.Response{error=response.ErrorBody}
//	@Failure		403		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Router			/api/v1/auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

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

	result, err := h.service.Login(r.Context(), req)
	if err != nil {
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
//	@Summary		Get current user
//	@Description	Returns the authenticated user's profile info
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=UserInfo}
//	@Failure		401	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/auth/me [get]
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
//
//	@Summary		Refresh token
//	@Description	Exchange a valid refresh token for new access and refresh tokens
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		RefreshRequest	true	"Refresh token"
//	@Success		200		{object}	response.Response{data=TokenResponse}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		401		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Router			/api/v1/auth/refresh [post]
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
//	@Summary		Logout
//	@Description	Logout current user (client should discard tokens)
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	response.Response{data=map[string]string}
//	@Security		BearerAuth
//	@Router			/api/v1/auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

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
