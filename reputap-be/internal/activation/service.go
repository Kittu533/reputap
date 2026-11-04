package activation

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	repo            *Repository
	redirectBaseURL string
}

func NewService(repo *Repository, redirectBaseURL string) *Service {
	return &Service{
		repo:            repo,
		redirectBaseURL: redirectBaseURL,
	}
}

func (s *Service) Check(ctx context.Context, code string) (*CheckResponse, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return &CheckResponse{
			Status:      "not_found",
			CanActivate: false,
			Message:     "Kode kartu wajib diisi",
		}, nil
	}

	qr, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if err == ErrNotFound {
			return &CheckResponse{
				Code:        code,
				Status:      "not_found",
				CanActivate: false,
				Message:     "Kartu tidak ditemukan di sistem REPUTAP",
			}, nil
		}
		return nil, err
	}

	if qr.Status == "active" && qr.StoreID != nil && *qr.StoreID > 0 {
		return &CheckResponse{
			Code:        code,
			Status:      "active",
			CanActivate: false,
			Message:     fmt.Sprintf("Kartu ini sudah aktif dan terhubung ke toko '%s'", qr.StoreName),
			StoreName:   qr.StoreName,
			ReviewURL:   qr.ReviewURL,
		}, nil
	}

	return &CheckResponse{
		Code:        code,
		Status:      "unassigned",
		CanActivate: true,
		Message:     "Kartu siap diaktivasi",
	}, nil
}

func (s *Service) Activate(ctx context.Context, req ActivateRequest) (*ActivateResponse, error) {
	req.Code = strings.TrimSpace(req.Code)
	req.BusinessName = strings.TrimSpace(req.BusinessName)
	req.ReviewURL = strings.TrimSpace(req.ReviewURL)

	_, err := s.repo.Activate(ctx, req)
	if err != nil {
		return nil, err
	}

	return &ActivateResponse{
		Code:        req.Code,
		StoreName:   req.BusinessName,
		ReviewURL:   req.ReviewURL,
		RedirectURL: fmt.Sprintf("%s/%s", s.redirectBaseURL, req.Code),
	}, nil
}
