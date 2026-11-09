package qrcodes

import (
	"encoding/json"
	"errors"
	"fmt"
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
//	@Summary		List QR codes
//	@Description	Get paginated list of QR codes with optional filters
//	@Tags			QR Codes
//	@Produce		json
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			per_page	query		int		false	"Items per page"	default(20)
//	@Param			store_id	query		int		false	"Filter by store ID"
//	@Param			status		query		string	false	"Filter by status (active, inactive, suspended, retired)"
//	@Param			search		query		string	false	"Search by code"
//	@Success		200	{object}	response.Response{data=[]QRCodeDetail,meta=response.Meta}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/qr-codes [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pg := pagination.Parse(r)
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")
	storeID, _ := strconv.ParseUint(r.URL.Query().Get("store_id"), 10, 64)

	params := ListParams{
		Page:    pg.Page,
		PerPage: pg.PerPage,
		Offset:  pg.Offset,
		StoreID: storeID,
		Status:  status,
		Search:  search,
	}

	qrCodes, total, err := h.service.List(r.Context(), params)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Paginated(w, http.StatusOK, qrCodes, response.Meta{
		Page:       pg.Page,
		PerPage:    pg.PerPage,
		Total:      total,
		TotalPages: pagination.TotalPages(total, pg.PerPage),
	})
}

// GetByID godoc
//
//	@Summary		Get QR code by ID
//	@Description	Get a single QR code by its ID
//	@Tags			QR Codes
//	@Produce		json
//	@Param			id	path		int	true	"QR Code ID"
//	@Success		200	{object}	response.Response{data=QRCodeDetail}
//	@Failure		404	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/qr-codes/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid QR code ID")
		return
	}

	qr, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, qr)
}

// UpdateStatus godoc
//
//	@Summary		Update QR code status
//	@Description	Update the status of a QR code
//	@Tags			QR Codes
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"QR Code ID"
//	@Param			body	body		UpdateStatusRequest		true	"Status data"
//	@Success		200		{object}	response.Response{data=map[string]string}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		404		{object}	response.Response{error=response.ErrorBody}
//	@Failure		422		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/qr-codes/{id}/status [patch]
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid QR code ID")
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
		"message": "QR code status updated successfully",
	})
}

// Update godoc
//
//	@Summary		Update QR code
//	@Description	Update QR code status and assigned store
//	@Tags			QR Codes
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"QR Code ID"
//	@Param			request	body		UpdateQRCodeRequest	true	"Update request"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		404		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/qr-codes/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid QR code ID")
		return
	}

	var req UpdateQRCodeRequest
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

	if err := h.service.Update(r.Context(), id, req); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "QR code updated successfully",
	})
}

// Delete godoc
//
//	@Summary		Delete QR code
//	@Description	Delete a QR code and its associated product instance
//	@Tags			QR Codes
//	@Produce		json
//	@Param			id	path		int	true	"QR Code ID"
//	@Success		200	{object}	response.Response
//	@Failure		404	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/qr-codes/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid QR code ID")
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
		"message": "QR code deleted successfully",
	})
}

// BatchGenerate godoc
//
//	@Summary		Batch generate QR codes
//	@Description	Generate a batch of blank/unassigned QR codes directly from dashboard
//	@Tags			QR Codes
//	@Accept			json
//	@Produce		json
//	@Param			request	body		BatchGenerateRequest	true	"Batch generate request"
//	@Success		201		{object}	response.Response{data=[]QRCodeDetail}
//	@Failure		400		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/qr-codes/batch [post]
func (h *Handler) BatchGenerate(w http.ResponseWriter, r *http.Request) {
	var req BatchGenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Prefix == "" {
		req.Prefix = "RP26"
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

	results, err := h.service.BatchGenerate(r.Context(), req)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, results)
}

// Download godoc
//
//	@Summary		Download QR code image
//	@Description	Download QR code as a PNG image
//	@Tags			QR Codes
//	@Produce		png
//	@Param			id		path		int	true	"QR Code ID"
//	@Param			size	query		int	false	"Image size in pixels"	default(256)
//	@Success		200		{file}		binary	"QR code PNG image"
//	@Failure		404		{object}	response.Response{error=response.ErrorBody}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/qr-codes/{id}/download [get]
func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_ID", "Invalid QR code ID")
		return
	}

	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	qr, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(w, err.Error())
			return
		}
		response.InternalError(w)
		return
	}

	pngBytes, err := h.service.GenerateQRImage(r.Context(), id, size)
	if err != nil {
		response.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.png", qr.Code))
	w.WriteHeader(http.StatusOK)
	w.Write(pngBytes)
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
