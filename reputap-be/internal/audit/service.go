package audit

import (
	"context"
	"encoding/json"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Log(ctx context.Context, userID uint64, action string, entityType string, entityID *uint64, oldValues interface{}, newValues interface{}) error {
	var oldJSON *json.RawMessage
	var newJSON *json.RawMessage

	if oldValues != nil {
		data, err := json.Marshal(oldValues)
		if err != nil {
			return fmt.Errorf("marshaling old values: %w", err)
		}
		raw := json.RawMessage(data)
		oldJSON = &raw
	}

	if newValues != nil {
		data, err := json.Marshal(newValues)
		if err != nil {
			return fmt.Errorf("marshaling new values: %w", err)
		}
		raw := json.RawMessage(data)
		newJSON = &raw
	}

	log := &AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		OldValues:  oldJSON,
		NewValues:  newJSON,
	}

	if err := s.repo.Create(ctx, log); err != nil {
		return fmt.Errorf("creating audit log: %w", err)
	}

	return nil
}

func (s *Service) List(ctx context.Context, params ListParams) ([]AuditLog, int, error) {
	return s.repo.List(ctx, params)
}
