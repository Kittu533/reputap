package qrcodes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

var ErrNotFound = errors.New("qr code not found")

type Service struct {
	repo            *Repository
	redirectBaseURL string
}

func NewService(repo *Repository, redirectBaseURL string) *Service {
	return &Service{repo: repo, redirectBaseURL: redirectBaseURL}
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*QRCodeDetail, error) {
	qr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching qr code: %w", err)
	}

	detail := &QRCodeDetail{
		QRCode:      *qr,
		RedirectURL: s.buildRedirectURL(qr.Code),
	}

	return detail, nil
}

func (s *Service) List(ctx context.Context, params ListParams) ([]QRCodeDetail, int, error) {
	qrCodes, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("listing qr codes: %w", err)
	}

	details := make([]QRCodeDetail, len(qrCodes))
	for i, qr := range qrCodes {
		details[i] = QRCodeDetail{
			QRCode:      qr,
			RedirectURL: s.buildRedirectURL(qr.Code),
		}
	}

	return details, total, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id uint64, req UpdateStatusRequest) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("fetching qr code: %w", err)
	}

	if err := s.repo.UpdateStatus(ctx, id, req.Status); err != nil {
		return fmt.Errorf("updating qr code status: %w", err)
	}

	return nil
}

func (s *Service) GenerateQRImage(ctx context.Context, id uint64, size int) ([]byte, error) {
	qr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetching qr code: %w", err)
	}

	if size <= 0 {
		size = 256
	}

	url := s.buildRedirectURL(qr.Code)

	pngBytes, err := qrcode.Encode(url, qrcode.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("generating qr image: %w", err)
	}

	return pngBytes, nil
}

func (s *Service) buildRedirectURL(code string) string {
	return s.redirectBaseURL + "/" + code
}
