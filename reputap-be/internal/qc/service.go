package qc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("quality control record not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateRequest, userID uint64) (*QCWithInstance, error) {
	qc := &QualityControl{
		ProductInstanceID: req.ProductInstanceID,
		InspectedBy:       userID,
		Result:            req.Result,
		ChecklistData:     req.ChecklistData,
		EvidencePaths:     req.EvidencePaths,
		Notes:             req.Notes,
	}

	if err := s.repo.Create(ctx, qc); err != nil {
		return nil, fmt.Errorf("creating quality control: %w", err)
	}

	statusMap := map[string]string{
		"passed":          "passed",
		"failed":          "failed",
		"rework_required": "rework",
	}

	if status, ok := statusMap[req.Result]; ok {
		if err := s.repo.UpdateInstanceQCStatus(ctx, req.ProductInstanceID, status); err != nil {
			return nil, fmt.Errorf("updating instance qc status: %w", err)
		}
	}

	created, err := s.repo.GetByID(ctx, qc.ID)
	if err != nil {
		return nil, fmt.Errorf("fetching created quality control: %w", err)
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*QCWithInstance, error) {
	qc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching quality control: %w", err)
	}

	return qc, nil
}

func (s *Service) ListByInstanceID(ctx context.Context, instanceID uint64) ([]QualityControl, error) {
	return s.repo.ListByInstanceID(ctx, instanceID)
}

func (s *Service) List(ctx context.Context, params ListParams) ([]QCWithInstance, int, error) {
	return s.repo.List(ctx, params)
}
