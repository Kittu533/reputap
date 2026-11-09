package stores

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/kittu533/reputap-be/internal/response"
	"github.com/kittu533/reputap-be/pkg/pagination"
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

// List godoc
//
//	@Summary		List stores
//	@Description	Get paginated list of stores with optional filters
//	@Tags			Stores
//	@Produce		json
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			per_page	query		int		false	"Items per page"	default(20)
//	@Param			search		query		string	false	"Search by name or store code"
//	@Param			customer_id	query		int		false	"Filter by customer ID"
//	@Param			status		query		string	false	"Filter by status (active, inactive, suspended)"
//	@Success		200	{object}	response.Response{data=[]Store,meta=response.Meta}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/stores [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pg := pagination.Parse(r)
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")
	customerID, _ := strconv.ParseUint(r.URL.Query().Get("customer_id"), 10, 64)

	params := ListParams{
		Page:       pg.Page,
		PerPage:    pg.PerPage,
		Offset:     pg.Offset,
		Search:     search,
		CustomerID: customerID,
		Status:     status,
	}

	stores, total, err := h.service.List(r.Context(), params)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Paginated(w, http.StatusOK, stores, response.Meta{
		Page:       pg.Page,
		PerPage:    pg.PerPage,
		Total:      total,
		TotalPages: pagination.TotalPages(total, pg.PerPage),
	})
}

// GetByID godoc
//
//	@Summary		Get store by ID
//	@Description	Get a single store by its ID
//	@Tags			Stores
//	@Produce		json
//	@Param			id	path		int	true	"Store ID"
//	@Success		200	{object}	response.Response{data=Store}
//	@Failure		404	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/stores/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid store ID")
		return
	}

	store, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, store)
}

// Create godoc
//
//	@Summary		Create store
//	@Description	Create a new store with auto-generated store code and slug
//	@Tags			Stores
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateRequest	true	"Store data"
//	@Success		201		{object}	response.Response{data=Store}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		422		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/stores [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest

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

	store, err := h.service.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidReviewURL) {
			response.Err(w, http.StatusBadRequest, "INVALID_REVIEW_URL", err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, store)
}

// Update godoc
//
//	@Summary		Update store
//	@Description	Update an existing store by ID
//	@Tags			Stores
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Store ID"
//	@Param			body	body		UpdateRequest	true	"Store data"
//	@Success		200		{object}	response.Response{data=Store}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		404		{object}	response.Response{error=response.ErrorBody}
//	@Failure		422		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/stores/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid store ID")
		return
	}

	var req UpdateRequest

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

	store, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, store)
}

// UpdateStatus godoc
//
//	@Summary		Update store status
//	@Description	Update the status of a store
//	@Tags			Stores
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Store ID"
//	@Param			body	body		UpdateStatusRequest	true	"Status data"
//	@Success		200		{object}	response.Response{data=map[string]string}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		404		{object}	response.Response{error=response.ErrorBody}
//	@Failure		422		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/stores/{id}/status [patch]
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid store ID")
		return
	}

	var req UpdateStatusRequest

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

	if err := h.service.UpdateStatus(r.Context(), id, req); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Store status updated successfully",
	})
}

// UpdateReviewURL godoc
//
//	@Summary		Update store review URL
//	@Description	Update the Google Review URL of a store
//	@Tags			Stores
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Store ID"
//	@Param			body	body		UpdateReviewURLRequest	true	"Review URL data"
//	@Success		200		{object}	response.Response{data=map[string]string}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		404		{object}	response.Response{error=response.ErrorBody}
//	@Failure		422		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/stores/{id}/review-url [put]
func (h *Handler) UpdateReviewURL(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid store ID")
		return
	}

	var req UpdateReviewURLRequest

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

	if err := h.service.UpdateReviewURL(r.Context(), id, req); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		if errors.Is(err, ErrInvalidReviewURL) {
			response.Err(w, http.StatusBadRequest, "INVALID_REVIEW_URL", err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Store review URL updated successfully",
	})
}

// Delete godoc
//
//	@Summary		Delete store
//	@Description	Delete a store by ID
//	@Tags			Stores
//	@Produce		json
//	@Param			id	path		int	true	"Store ID"
//	@Success		200	{object}	response.Response
//	@Failure		404	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/stores/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid store ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.Err(w, http.StatusBadRequest, "DELETE_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Store deleted successfully",
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
	case "oneof":
		return e.Field() + " must be one of: " + e.Param()
	default:
		return e.Field() + " is invalid"
	}
}
