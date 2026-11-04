package activation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/kittu533/reputap-be/internal/response"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service:  service,
		validate: validator.New(),
	}
}

// Check godoc
//
//	@Summary		Check card activation status
//	@Description	Check whether a QR/NFC card code exists and can be activated
//	@Tags			Activation
//	@Produce		json
//	@Param			code	query		string	true	"Card code"
//	@Success		200		{object}	response.Response{data=CheckResponse}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Router			/api/v1/activate/check [get]
func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	res, err := h.service.Check(r.Context(), code)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

// Activate godoc
//
//	@Summary		Public self-service card activation
//	@Description	Activate a blank / unassigned card and bind it to a store
//	@Tags			Activation
//	@Accept			json
//	@Produce		json
//	@Param			body	body		ActivateRequest	true	"Activation data"
//	@Success		200		{object}	response.Response{data=ActivateResponse}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		404		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Router			/api/v1/activate [post]
func (h *Handler) Activate(w http.ResponseWriter, r *http.Request) {
	var req ActivateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Format data tidak valid")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var details []response.FieldError
			for _, e := range validationErrors {
				details = append(details, response.FieldError{
					Field:   e.Field(),
					Message: e.Error(),
				})
			}
			response.ValidationErr(w, details)
			return
		}
		response.Err(w, http.StatusBadRequest, "VALIDATION_ERROR", "Data form tidak lengkap")
		return
	}

	res, err := h.service.Activate(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, "Kode kartu tidak terdaftar di sistem")
			return
		}
		if errors.Is(err, ErrAlreadyActive) {
			response.Err(w, http.StatusBadRequest, "ALREADY_ACTIVE", "Kartu ini sudah diaktivasi sebelumnya")
			return
		}
		response.Err(w, http.StatusInternalServerError, "ACTIVATION_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, res)
}
