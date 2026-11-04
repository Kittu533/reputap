package production

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

// Create godoc
//
//	@Summary		Create production job
//	@Description	Create a new production job for a product instance
//	@Tags			Production
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateRequest	true	"Production job data"
//	@Success		201		{object}	response.Response{data=JobWithInstance}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		422		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/production/jobs [post]
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

	job, err := h.service.Create(r.Context(), req)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, job)
}

// List godoc
//
//	@Summary		List production jobs
//	@Description	Get paginated list of production jobs with optional filters
//	@Tags			Production
//	@Produce		json
//	@Param			page		query		int		false	"Page number"			default(1)
//	@Param			per_page	query		int		false	"Items per page"		default(20)
//	@Param			status		query		string	false	"Filter by status"
//	@Param			assigned_to	query		int		false	"Filter by assigned user ID"
//	@Param			store_id	query		int		false	"Filter by store ID"
//	@Success		200			{object}	response.Response{data=[]JobWithInstance,meta=response.Meta}
//	@Failure		500			{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/production/jobs [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pg := pagination.Parse(r)

	assignedTo, _ := strconv.ParseUint(r.URL.Query().Get("assigned_to"), 10, 64)
	storeID, _ := strconv.ParseUint(r.URL.Query().Get("store_id"), 10, 64)

	params := ListParams{
		Page:       pg.Page,
		PerPage:    pg.PerPage,
		Offset:     pg.Offset,
		Status:     r.URL.Query().Get("status"),
		AssignedTo: assignedTo,
		StoreID:    storeID,
	}

	jobs, total, err := h.service.List(r.Context(), params)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Paginated(w, http.StatusOK, jobs, response.Meta{
		Page:       pg.Page,
		PerPage:    pg.PerPage,
		Total:      total,
		TotalPages: pagination.TotalPages(total, pg.PerPage),
	})
}

// GetByID godoc
//
//	@Summary		Get production job by ID
//	@Description	Get a single production job by ID
//	@Tags			Production
//	@Produce		json
//	@Param			id	path		int	true	"Production Job ID"
//	@Success		200	{object}	response.Response{data=JobWithInstance}
//	@Failure		404	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/production/jobs/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid production job ID")
		return
	}

	job, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, job)
}

// UpdateStatus godoc
//
//	@Summary		Update production job status
//	@Description	Update the status of a production job
//	@Tags			Production
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Production Job ID"
//	@Param			body	body		UpdateStatusRequest	true	"Status data"
//	@Success		200		{object}	response.Response{data=map[string]string}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		404		{object}	response.Response{error=response.ErrorBody}
//	@Failure		422		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/production/jobs/{id}/status [patch]
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid production job ID")
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
		"message": "Production job status updated successfully",
	})
}

func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "email":
		return e.Field() + " must be a valid email"
	case "min":
		return e.Field() + " must be at least " + e.Param()
	case "max":
		return e.Field() + " must be at most " + e.Param()
	case "gt":
		return e.Field() + " must be greater than " + e.Param()
	case "gte":
		return e.Field() + " must be greater than or equal to " + e.Param()
	case "oneof":
		return e.Field() + " must be one of: " + e.Param()
	default:
		return e.Field() + " is invalid"
	}
}
