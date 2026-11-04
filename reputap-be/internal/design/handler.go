package design

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/kittu533/reputap-be/internal/auth"
	"github.com/kittu533/reputap-be/internal/response"
	"github.com/kittu533/reputap-be/pkg/pagination"
)

type Handler struct {
	service       *Service
	validate      *validator.Validate
	maxUploadSize int64
}

func NewHandler(service *Service, maxUploadSize int64) *Handler {
	return &Handler{
		service:       service,
		validate:      validator.New(),
		maxUploadSize: maxUploadSize,
	}
}

// Upload godoc
//
//	@Summary		Upload design asset
//	@Description	Upload a design file for an order
//	@Tags			Design Assets
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file		formData	file	true	"Design file"
//	@Param			order_id	formData	int		true	"Order ID"
//	@Param			store_id	formData	int		true	"Store ID"
//	@Success		201	{object}	response.Response{data=DesignAsset}
//	@Failure		400	{object}	response.Response{error=response.ErrorBody}
//	@Failure		401	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/design-assets/upload [post]
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "Not authenticated")
		return
	}

	if err := r.ParseMultipartForm(h.maxUploadSize); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_FORM", "Failed to parse multipart form")
		return
	}

	orderID, err := strconv.ParseUint(r.FormValue("order_id"), 10, 64)
	if err != nil || orderID == 0 {
		response.Err(w, http.StatusBadRequest, "INVALID_ORDER_ID", "Invalid or missing order_id")
		return
	}

	storeID, err := strconv.ParseUint(r.FormValue("store_id"), 10, 64)
	if err != nil || storeID == 0 {
		response.Err(w, http.StatusBadRequest, "INVALID_STORE_ID", "Invalid or missing store_id")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Err(w, http.StatusBadRequest, "MISSING_FILE", "File is required")
		return
	}
	defer file.Close()

	asset, err := h.service.Upload(r.Context(), orderID, storeID, file, header, claims.UserID)
	if err != nil {
		if errors.Is(err, ErrInvalidFileType) {
			response.Err(w, http.StatusBadRequest, "INVALID_FILE_TYPE", err.Error())
			return
		}
		if errors.Is(err, ErrFileTooLarge) {
			response.Err(w, http.StatusBadRequest, "FILE_TOO_LARGE", err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, asset)
}

// List godoc
//
//	@Summary		List design assets
//	@Description	Get paginated list of design assets with optional filters
//	@Tags			Design Assets
//	@Produce		json
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			per_page	query		int		false	"Items per page"	default(20)
//	@Param			order_id	query		int		false	"Filter by order ID"
//	@Param			store_id	query		int		false	"Filter by store ID"
//	@Param			status		query		string	false	"Filter by approval status"
//	@Success		200	{object}	response.Response{data=[]DesignAsset,meta=response.Meta}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/design-assets [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pg := pagination.Parse(r)

	orderID, _ := strconv.ParseUint(r.URL.Query().Get("order_id"), 10, 64)
	storeID, _ := strconv.ParseUint(r.URL.Query().Get("store_id"), 10, 64)

	params := ListParams{
		Page:    pg.Page,
		PerPage: pg.PerPage,
		Offset:  pg.Offset,
		OrderID: orderID,
		StoreID: storeID,
		Status:  r.URL.Query().Get("status"),
	}

	assets, total, err := h.service.List(r.Context(), params)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Paginated(w, http.StatusOK, assets, response.Meta{
		Page:       pg.Page,
		PerPage:    pg.PerPage,
		Total:      total,
		TotalPages: pagination.TotalPages(total, pg.PerPage),
	})
}

// ListByOrderID godoc
//
//	@Summary		List designs by order
//	@Description	Get all design assets for a specific order
//	@Tags			Design Assets
//	@Produce		json
//	@Param			orderId	path		int	true	"Order ID"
//	@Success		200	{object}	response.Response{data=[]DesignAsset}
//	@Failure		400	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/orders/{orderId}/designs [get]
func (h *Handler) ListByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseUint(chi.URLParam(r, "orderId"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid order ID")
		return
	}

	assets, err := h.service.ListByOrderID(r.Context(), orderID)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, assets)
}

// GetByID godoc
//
//	@Summary		Get design asset by ID
//	@Description	Get a single design asset by its ID
//	@Tags			Design Assets
//	@Produce		json
//	@Param			id	path		int	true	"Design Asset ID"
//	@Success		200	{object}	response.Response{data=DesignAsset}
//	@Failure		400	{object}	response.Response{error=response.ErrorBody}
//	@Failure		404	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/design-assets/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid design asset ID")
		return
	}

	asset, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, asset)
}

// UpdateApprovalStatus godoc
//
//	@Summary		Update approval status
//	@Description	Update the approval status of a design asset
//	@Tags			Design Assets
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Design Asset ID"
//	@Param			body	body		UpdateApprovalRequest	true	"Approval status data"
//	@Success		200	{object}	response.Response{data=map[string]string}
//	@Failure		400	{object}	response.Response{error=response.ErrorBody}
//	@Failure		404	{object}	response.Response{error=response.ErrorBody}
//	@Failure		422	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/design-assets/{id}/approval [patch]
func (h *Handler) UpdateApprovalStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid design asset ID")
		return
	}

	var req UpdateApprovalRequest

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

	if err := h.service.UpdateApprovalStatus(r.Context(), id, req); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Approval status updated successfully",
	})
}

// Delete godoc
//
//	@Summary		Delete design asset
//	@Description	Delete a design asset and its file from disk
//	@Tags			Design Assets
//	@Produce		json
//	@Param			id	path		int	true	"Design Asset ID"
//	@Success		200	{object}	response.Response{data=map[string]string}
//	@Failure		400	{object}	response.Response{error=response.ErrorBody}
//	@Failure		404	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/design-assets/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid design asset ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Design asset deleted successfully",
	})
}

// Download godoc
//
//	@Summary		Download design asset
//	@Description	Download a design asset file
//	@Tags			Design Assets
//	@Produce		octet-stream
//	@Param			id	path		int	true	"Design Asset ID"
//	@Success		200	{file}		file
//	@Failure		400	{object}	response.Response{error=response.ErrorBody}
//	@Failure		404	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/design-assets/{id}/download [get]
func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid design asset ID")
		return
	}

	asset, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	contentType := detectContentType(asset.FileName)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", asset.FileName))
	http.ServeFile(w, r, asset.FilePath)
}

func detectContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
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
