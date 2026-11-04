package productinstances

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

// GenerateForOrder godoc
//
//	@Summary		Generate product instances for an order
//	@Description	Create product instance records and QR codes for all order items
//	@Tags			Product Instances
//	@Produce		json
//	@Param			orderId	path		int	true	"Order ID"
//	@Success		201		{object}	response.Response{data=[]InstanceWithQR}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		409		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/orders/{orderId}/instances/generate [post]
func (h *Handler) GenerateForOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseUint(chi.URLParam(r, "orderId"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid order ID")
		return
	}

	instances, err := h.service.GenerateForOrder(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, ErrAlreadyGenerated) {
			response.Err(w, http.StatusConflict, "ALREADY_GENERATED", err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, instances)
}

// ListByOrderID godoc
//
//	@Summary		List instances by order
//	@Description	Get all product instances for a specific order
//	@Tags			Product Instances
//	@Produce		json
//	@Param			orderId	path		int	true	"Order ID"
//	@Success		200		{object}	response.Response{data=[]InstanceWithQR}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/orders/{orderId}/instances [get]
func (h *Handler) ListByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseUint(chi.URLParam(r, "orderId"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid order ID")
		return
	}

	instances, err := h.service.ListByOrderID(r.Context(), orderID)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, instances)
}

// GetByID godoc
//
//	@Summary		Get product instance by ID
//	@Description	Get a single product instance with its QR code
//	@Tags			Product Instances
//	@Produce		json
//	@Param			id	path		int	true	"Product Instance ID"
//	@Success		200	{object}	response.Response{data=InstanceWithQR}
//	@Failure		400	{object}	response.Response{error=response.ErrorBody}
//	@Failure		404	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/product-instances/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid product instance ID")
		return
	}

	instance, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, instance)
}

// List godoc
//
//	@Summary		List product instances
//	@Description	Get paginated list of product instances with optional filters
//	@Tags			Product Instances
//	@Produce		json
//	@Param			page				query		int		false	"Page number"				default(1)
//	@Param			per_page			query		int		false	"Items per page"			default(20)
//	@Param			store_id			query		int		false	"Filter by store ID"
//	@Param			nfc_status			query		string	false	"Filter by NFC status"
//	@Param			production_status	query		string	false	"Filter by production status"
//	@Success		200					{object}	response.Response{data=[]InstanceWithQR,meta=response.Meta}
//	@Failure		500					{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/product-instances [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pg := pagination.Parse(r)

	storeID, _ := strconv.ParseUint(r.URL.Query().Get("store_id"), 10, 64)

	params := ListParams{
		Page:       pg.Page,
		PerPage:    pg.PerPage,
		Offset:     pg.Offset,
		StoreID:    storeID,
		NFCStatus:  r.URL.Query().Get("nfc_status"),
		ProdStatus: r.URL.Query().Get("production_status"),
	}

	instances, total, err := h.service.List(r.Context(), params)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Paginated(w, http.StatusOK, instances, response.Meta{
		Page:       pg.Page,
		PerPage:    pg.PerPage,
		Total:      total,
		TotalPages: pagination.TotalPages(total, pg.PerPage),
	})
}

// UpdateNFCStatus godoc
//
//	@Summary		Update NFC status
//	@Description	Update the NFC configuration status of a product instance
//	@Tags			Product Instances
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Product Instance ID"
//	@Param			body	body		UpdateNFCStatusRequest	true	"NFC status data"
//	@Success		200		{object}	response.Response{data=map[string]string}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		401		{object}	response.Response{error=response.ErrorBody}
//	@Failure		404		{object}	response.Response{error=response.ErrorBody}
//	@Failure		422		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/product-instances/{id}/nfc-status [patch]
func (h *Handler) UpdateNFCStatus(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "Not authenticated")
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid product instance ID")
		return
	}

	var req UpdateNFCStatusRequest

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

	if err := h.service.UpdateNFCStatus(r.Context(), id, req, claims.UserID); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "NFC status updated successfully",
	})
}

func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "oneof":
		return e.Field() + " must be one of: " + e.Param()
	default:
		return e.Field() + " is invalid"
	}
}
