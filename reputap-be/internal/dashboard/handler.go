package dashboard

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/kittu533/reputap-be/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Summary godoc
//
//	@Summary		Get dashboard summary
//	@Description	Get key operational and business metrics summary
//	@Tags			Dashboard
//	@Produce		json
//	@Success		200	{object}	response.Response{data=SummaryResponse}
//	@Failure		500	{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/dashboard/summary [get]
func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.service.GetSummary(r.Context())
	if err != nil {
		slog.Error("failed to get dashboard summary", "error", err)
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, summary)
}

// RecentActivities godoc
//
//	@Summary		Get recent activities
//	@Description	Get latest audit log activities for dashboard feed
//	@Tags			Dashboard
//	@Produce		json
//	@Param			limit	query		int	false	"Number of activities"	default(10)
//	@Success		200		{object}	response.Response{data=[]RecentActivity}
//	@Failure		500		{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/dashboard/recent-activities [get]
func (h *Handler) RecentActivities(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	activities, err := h.service.GetRecentActivities(r.Context(), limit)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, activities)
}
