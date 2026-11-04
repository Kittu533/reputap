package production

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("production job not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*JobWithInstance, error) {
	job := &Job{
		ProductInstanceID: req.ProductInstanceID,
		AssignedTo:        req.AssignedTo,
		Status:            "pending",
		ChecklistData:     req.ChecklistData,
		Notes:             req.Notes,
	}

	if err := s.repo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("creating production job: %w", err)
	}

	created, err := s.repo.GetByID(ctx, job.ID)
	if err != nil {
		return nil, fmt.Errorf("fetching created production job: %w", err)
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*JobWithInstance, error) {
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching production job: %w", err)
	}

	return job, nil
}

func (s *Service) List(ctx context.Context, params ListParams) ([]JobWithInstance, int, error) {
	return s.repo.List(ctx, params)
}

func (s *Service) UpdateStatus(ctx context.Context, id uint64, req UpdateStatusRequest) error {
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("fetching production job: %w", err)
	}

	if err := s.repo.UpdateStatus(ctx, id, req.Status, req.Notes, req.ChecklistData); err != nil {
		return fmt.Errorf("updating production job status: %w", err)
	}

	if err := s.repo.UpdateProductInstanceStatus(ctx, job.ProductInstanceID, req.Status); err != nil {
		return fmt.Errorf("syncing product instance status: %w", err)
	}

	return nil
}
