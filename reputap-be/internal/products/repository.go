package products

import (
	"context"
	"strings"

	"github.com/kittu533/reputap-be/internal/database"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, product *Product) error {
	query := `
		INSERT INTO products (sku, name, description, material, width_mm, height_mm, thickness_mm, nfc_chip_type, finishing, production_cost, selling_price, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	id, err := r.db.InsertID(ctx, query,
		product.SKU,
		product.Name,
		product.Description,
		product.Material,
		product.WidthMM,
		product.HeightMM,
		product.ThicknessMM,
		product.NFCChipType,
		product.Finishing,
		product.ProductionCost,
		product.SellingPrice,
		product.Status,
	)
	if err != nil {
		return err
	}

	product.ID = id
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*Product, error) {
	var product Product

	query := `SELECT * FROM products WHERE id = ?`

	err := r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]Product, int, error) {
	var total int
	var products []Product

	baseWhere := ""
	var conditions []string
	var args []interface{}

	if strings.TrimSpace(params.Search) != "" {
		pattern := "%" + params.Search + "%"
		conditions = append(conditions, "(name LIKE ? OR sku LIKE ?)")
		args = append(args, pattern, pattern)
	}

	if strings.TrimSpace(params.Status) != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, params.Status)
	}

	if len(conditions) > 0 {
		baseWhere = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM products" + baseWhere
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := "SELECT * FROM products" + baseWhere + " ORDER BY id DESC LIMIT ? OFFSET ?"
	selectArgs := append(args, params.PerPage, params.Offset)

	err = r.db.SelectContext(ctx, &products, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *Repository) Update(ctx context.Context, product *Product) error {
	query := `
		UPDATE products
		SET sku = ?, name = ?, description = ?, material = ?, width_mm = ?, height_mm = ?, thickness_mm = ?,
		    nfc_chip_type = ?, finishing = ?, production_cost = ?, selling_price = ?, status = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		product.SKU,
		product.Name,
		product.Description,
		product.Material,
		product.WidthMM,
		product.HeightMM,
		product.ThicknessMM,
		product.NFCChipType,
		product.Finishing,
		product.ProductionCost,
		product.SellingPrice,
		product.Status,
		product.ID,
	)
	return err
}

func (r *Repository) SKUExists(ctx context.Context, sku string, excludeID uint64) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM products WHERE sku = ? AND id != ?`
	err := r.db.GetContext(ctx, &count, query, sku, excludeID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
