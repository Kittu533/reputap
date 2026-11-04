package activation

type CheckResponse struct {
	Code        string `json:"code"`
	Status      string `json:"status"` // "unassigned", "active", "not_found"
	CanActivate bool   `json:"can_activate"`
	Message     string `json:"message"`
	StoreName   string `json:"store_name,omitempty"`
	ReviewURL   string `json:"review_url,omitempty"`
}

type ActivateRequest struct {
	Code         string `json:"code" validate:"required"`
	BusinessName string `json:"business_name" validate:"required,min=2,max=150"`
	Category     string `json:"category" validate:"required"`
	ReviewURL    string `json:"review_url" validate:"required,url"`
	OwnerName    string `json:"owner_name" validate:"required,min=2,max=100"`
	Phone        string `json:"phone" validate:"required,min=8,max=20"`
	Email        string `json:"email" validate:"omitempty,email"`
	Address      string `json:"address"`
	City         string `json:"city"`
}

type ActivateResponse struct {
	Code        string `json:"code"`
	StoreName   string `json:"store_name"`
	ReviewURL   string `json:"review_url"`
	RedirectURL string `json:"redirect_url"`
}
