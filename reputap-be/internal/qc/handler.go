package qc

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/kittu533/reputap-be/internal/auth"
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
//	@Summary		Create quality control record
//	@Description	Create a new quality control inspection record
//	@Tags			Quality Control
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateRequest	true	"Quality control data"
//	@Success		201		{object}	response.Response{data=QCWithInstance}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		401		{object}	response.Response{error=response.ErrorBody}
//	@Failure		422		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/quality-controls [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "Not authenticated")
		return
	}

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

	qc, err := h.service.Create(r.Context(), req, claims.UserID)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, qc)
}

// GetByID godoc
//
//	@Summary		Get quality control by ID
//	@Description	Get a single quality control record by ID
//	@Tags			Quality Control
//	@Produce		json
//	@Param			id	path		int	true	"Quality Control ID"
//	@Success		200	{object}	response.Response{data=QCWithInstance}
//	@Failure		404	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/quality-controls/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid quality control ID")
		return
	}

	qc, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, qc)
}

// ListByInstanceID godoc
//
//	@Summary		List quality controls by instance
//	@Description	Get all quality control records for a product instance
//	@Tags			Quality Control
//	@Produce		json
//	@Param			instanceId	path		int	true	"Product Instance ID"
//	@Success		200			{object}	response.Response{data=[]QualityControl}
//	@Failure		500			{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/quality-controls/instance/{instanceId} [get]
func (h *Handler) ListByInstanceID(w http.ResponseWriter, r *http.Request) {
	instanceID, err := strconv.ParseUint(chi.URLParam(r, "instanceId"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid product instance ID")
		return
	}

	qcs, err := h.service.ListByInstanceID(r.Context(), instanceID)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, qcs)
}

// List godoc
//
//	@Summary		List quality controls
//	@Description	Get paginated list of quality control records with optional filters
//	@Tags			Quality Control
//	@Produce		json
//	@Param			page				query		int		false	"Page number"			default(1)
//	@Param			per_page			query		int		false	"Items per page"		default(20)
//	@Param			product_instance_id	query		int		false	"Filter by product instance ID"
//	@Param			result				query		string	false	"Filter by result"
//	@Success		200					{object}	response.Response{data=[]QCWithInstance,meta=response.Meta}
//	@Failure		500					{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/quality-controls [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pg := pagination.Parse(r)

	productInstanceID, _ := strconv.ParseUint(r.URL.Query().Get("product_instance_id"), 10, 64)

	params := ListParams{
		Page:              pg.Page,
		PerPage:           pg.PerPage,
		Offset:            pg.Offset,
		ProductInstanceID: productInstanceID,
		Result:            r.URL.Query().Get("result"),
	}

	qcs, total, err := h.service.List(r.Context(), params)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Paginated(w, http.StatusOK, qcs, response.Meta{
		Page:       pg.Page,
		PerPage:    pg.PerPage,
		Total:      total,
		TotalPages: pagination.TotalPages(total, pg.PerPage),
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
