package stores

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/kittu533/reputap-be/pkg/googleplaces"
	urlValidator "github.com/kittu533/reputap-be/pkg/validator"
)

var (
	ErrNotFound         = errors.New("store not found")
	ErrInvalidReviewURL = errors.New("review URL must be a valid Google Review URL (https://)")
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9-]+`)
var multipleHyphens = regexp.MustCompile(`-{2,}`)

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = nonAlphanumeric.ReplaceAllString(slug, "-")
	slug = multipleHyphens.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Store, error) {
	if req.ReviewURL != nil && *req.ReviewURL != "" {
		normalized := googleplaces.ToDirectReviewURL(ctx, *req.ReviewURL)
		req.ReviewURL = &normalized
		if !urlValidator.IsValidGoogleReviewURL(*req.ReviewURL) {
			return nil, ErrInvalidReviewURL
		}
	}

	code, err := s.repo.GetNextCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("generating store code: %w", err)
	}

	slug := generateSlug(req.Name)
	slug, err = s.ensureUniqueSlug(ctx, slug, 0)
	if err != nil {
		return nil, fmt.Errorf("generating slug: %w", err)
	}

	category := req.Category
	if category == "" {
		category = "other"
	}

	verified := false
	if req.ReviewURL != nil && *req.ReviewURL != "" {
		verified = urlValidator.IsValidGoogleReviewURL(*req.ReviewURL)
	}

	store := &Store{
		CustomerID:         req.CustomerID,
		StoreCode:          code,
		Name:               req.Name,
		Slug:               slug,
		Category:           category,
		Address:            req.Address,
		GoogleBusinessName: req.GoogleBusinessName,
		ReviewURL:          req.ReviewURL,
		ReviewURLVerified:  verified,
		Status:             "active",
	}

	if err := s.repo.Create(ctx, store); err != nil {
		return nil, fmt.Errorf("creating store: %w", err)
	}

	created, err := s.repo.GetByID(ctx, store.ID)
	if err != nil {
		return nil, fmt.Errorf("fetching created store: %w", err)
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*Store, error) {
	store, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching store: %w", err)
	}
	return store, nil
}

func (s *Service) List(ctx context.Context, params ListParams) ([]Store, int, error) {
	return s.repo.List(ctx, params)
}

func (s *Service) Update(ctx context.Context, id uint64, req UpdateRequest) (*Store, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching store: %w", err)
	}

	if existing.Name != req.Name {
		slug := generateSlug(req.Name)
		slug, err = s.ensureUniqueSlug(ctx, slug, id)
		if err != nil {
			return nil, fmt.Errorf("generating slug: %w", err)
		}
		existing.Slug = slug
	}

	existing.Name = req.Name
	existing.Category = req.Category
	existing.Address = req.Address
	existing.GoogleBusinessName = req.GoogleBusinessName

	if existing.Category == "" {
		existing.Category = "other"
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("updating store: %w", err)
	}

	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetching updated store: %w", err)
	}

	return updated, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id uint64, req UpdateStatusRequest) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("fetching store: %w", err)
	}

	if err := s.repo.UpdateStatus(ctx, id, req.Status); err != nil {
		return fmt.Errorf("updating store status: %w", err)
	}

	return nil
}

func (s *Service) UpdateReviewURL(ctx context.Context, id uint64, req UpdateReviewURLRequest) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("fetching store: %w", err)
	}

	req.ReviewURL = googleplaces.ToDirectReviewURL(ctx, req.ReviewURL)
	if !urlValidator.IsValidGoogleReviewURL(req.ReviewURL) {
		return ErrInvalidReviewURL
	}

	if err := s.repo.UpdateReviewURL(ctx, id, req.ReviewURL, true); err != nil {
		return fmt.Errorf("updating store review URL: %w", err)
	}

	return nil
}

func (s *Service) ensureUniqueSlug(ctx context.Context, slug string, excludeID uint64) (string, error) {
	exists, err := s.repo.SlugExists(ctx, slug, excludeID)
	if err != nil {
		return "", err
	}
	if !exists {
		return slug, nil
	}

	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s-%d", slug, i)
		exists, err := s.repo.SlugExists(ctx, candidate, excludeID)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}
}
