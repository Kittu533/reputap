package design

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrNotFound        = errors.New("design asset not found")
	ErrInvalidFileType = errors.New("invalid file type, allowed: png, jpg, jpeg, svg, pdf")
	ErrFileTooLarge    = errors.New("file too large")
)

var allowedTypes = map[string]bool{
	"image/png":       true,
	"image/jpeg":      true,
	"image/svg+xml":   true,
	"application/pdf": true,
}

var allowedExts = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".svg":  true,
	".pdf":  true,
}

type Service struct {
	repo          *Repository
	uploadDir     string
	maxUploadSize int64
}

func NewService(repo *Repository, uploadDir string, maxUploadSize int64) *Service {
	return &Service{
		repo:          repo,
		uploadDir:     uploadDir,
		maxUploadSize: maxUploadSize,
	}
}

func (s *Service) Upload(ctx context.Context, orderID uint64, storeID uint64, file multipart.File, header *multipart.FileHeader, userID uint64) (*DesignAsset, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExts[ext] {
		return nil, ErrInvalidFileType
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return nil, ErrInvalidFileType
	}

	if header.Size > s.maxUploadSize {
		return nil, ErrFileTooLarge
	}

	latestVersion, err := s.repo.GetLatestVersionByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("getting latest version: %w", err)
	}
	newVersion := latestVersion + 1

	dir := filepath.Join(s.uploadDir, "designs", fmt.Sprintf("%d", orderID))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating upload directory: %w", err)
	}

	uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	filePath := filepath.Join(dir, uniqueName)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("creating file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("saving file: %w", err)
	}

	asset := &DesignAsset{
		OrderID:        orderID,
		StoreID:        storeID,
		FilePath:       filePath,
		FileName:       header.Filename,
		FileType:       contentType,
		FileSize:       uint64(header.Size),
		Version:        uint(newVersion),
		ApprovalStatus: "draft",
		UploadedBy:     userID,
	}

	if err := s.repo.Create(ctx, asset); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("creating design asset record: %w", err)
	}

	created, err := s.repo.GetByID(ctx, asset.ID)
	if err != nil {
		return nil, fmt.Errorf("fetching created asset: %w", err)
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*DesignAsset, error) {
	asset, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching design asset: %w", err)
	}

	return asset, nil
}

func (s *Service) ListByOrderID(ctx context.Context, orderID uint64) ([]DesignAsset, error) {
	return s.repo.ListByOrderID(ctx, orderID)
}

func (s *Service) List(ctx context.Context, params ListParams) ([]DesignAsset, int, error) {
	return s.repo.List(ctx, params)
}

func (s *Service) UpdateApprovalStatus(ctx context.Context, id uint64, req UpdateApprovalRequest) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("fetching design asset: %w", err)
	}

	if err := s.repo.UpdateApprovalStatus(ctx, id, req.Status); err != nil {
		return fmt.Errorf("updating approval status: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	asset, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("fetching design asset: %w", err)
	}

	os.Remove(asset.FilePath)

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting design asset: %w", err)
	}

	return nil
}
