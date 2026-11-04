package analytics

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kittu533/reputap-be/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Overview godoc
//
//	@Summary		Get analytics overview
//	@Description	Get overall scan statistics, trends, device breakdowns, and top stores
//	@Tags			Analytics
//	@Produce		json
//	@Success		200	{object}	response.Response{data=OverviewResponse}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/analytics/overview [get]
func (h *Handler) Overview(w http.ResponseWriter, r *http.Request) {
	overview, err := h.service.GetOverview(r.Context())
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, overview)
}

// Scans godoc
//
//	@Summary		Get scan time-series
//	@Description	Get scan activity over time with daily/weekly/monthly aggregation and optional store filter
//	@Tags			Analytics
//	@Produce		json
//	@Param			period		query		string	false	"Aggregation period: daily, weekly, monthly"	default(daily)
//	@Param			from		query		string	false	"Start date (YYYY-MM-DD)"
//	@Param			to			query		string	false	"End date (YYYY-MM-DD)"
//	@Param			store_id	query		int		false	"Filter by Store ID"
//	@Success		200			{object}	response.Response{data=[]TimeSeriesPoint}
//	@Failure		500			{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/analytics/scans [get]
func (h *Handler) Scans(w http.ResponseWriter, r *http.Request) {
	var storeID uint64
	if s := r.URL.Query().Get("store_id"); s != "" {
		if id, err := strconv.ParseUint(s, 10, 64); err == nil {
			storeID = id
		}
	}

	params := ScansParams{
		Period:  r.URL.Query().Get("period"),
		From:    r.URL.Query().Get("from"),
		To:      r.URL.Query().Get("to"),
		StoreID: storeID,
	}

	points, err := h.service.GetScansTimeSeries(r.Context(), params)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, points)
}

// StoreAnalytics godoc
//
//	@Summary		Get store-specific analytics
//	@Description	Get scan metrics, active QR codes, device breakdown, and recent daily scans for a store
//	@Tags			Analytics
//	@Produce		json
//	@Param			id	path		int	true	"Store ID"
//	@Success		200	{object}	response.Response{data=StoreAnalyticsResponse}
//	@Failure		404	{object}	response.Response{error=response.ErrorBody}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/analytics/stores/{id} [get]
func (h *Handler) StoreAnalytics(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "invalid store id")
		return
	}

	res, err := h.service.GetStoreAnalytics(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrStoreNotFound) {
			response.NotFound(w, "store not found")
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, res)
}
