package audit

import (
	"net/http"
	"strconv"

	"github.com/kittu533/reputap-be/internal/response"
	"github.com/kittu533/reputap-be/pkg/pagination"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// List godoc
//
//	@Summary		List audit logs
//	@Description	Get paginated list of audit logs with optional filters
//	@Tags			Audit Logs
//	@Produce		json
//	@Param			page		query		int		false	"Page number"				default(1)
//	@Param			per_page	query		int		false	"Items per page"			default(20)
//	@Param			entity_type	query		string	false	"Filter by entity type"
//	@Param			entity_id	query		int		false	"Filter by entity ID"
//	@Param			user_id		query		int		false	"Filter by user ID"
//	@Param			action		query		string	false	"Filter by action"
//	@Param			date_from	query		string	false	"Filter from date (YYYY-MM-DD)"
//	@Param			date_to		query		string	false	"Filter to date (YYYY-MM-DD)"
//	@Success		200			{object}	response.Response{data=[]AuditLog,meta=response.Meta}
//	@Failure		500			{object}	response.Response{error=response.ErrorBody}
//	@Security		BearerAuth
//	@Router			/api/v1/audit-logs [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pg := pagination.Parse(r)

	entityID, _ := strconv.ParseUint(r.URL.Query().Get("entity_id"), 10, 64)
	userID, _ := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)

	params := ListParams{
		Page:       pg.Page,
		PerPage:    pg.PerPage,
		Offset:     pg.Offset,
		EntityType: r.URL.Query().Get("entity_type"),
		EntityID:   entityID,
		UserID:     userID,
		Action:     r.URL.Query().Get("action"),
		DateFrom:   r.URL.Query().Get("date_from"),
		DateTo:     r.URL.Query().Get("date_to"),
	}

	logs, total, err := h.service.List(r.Context(), params)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Paginated(w, http.StatusOK, logs, response.Meta{
		Page:       pg.Page,
		PerPage:    pg.PerPage,
		Total:      total,
		TotalPages: pagination.TotalPages(total, pg.PerPage),
	})
}
