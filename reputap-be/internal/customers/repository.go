package customers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, customer *Customer) error {
	query := `
		INSERT INTO customers (customer_code, name, business_name, email, phone, address, city, province, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		customer.CustomerCode,
		customer.Name,
		customer.BusinessName,
		customer.Email,
		customer.Phone,
		customer.Address,
		customer.City,
		customer.Province,
		customer.Notes,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	customer.ID = uint64(id)
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*Customer, error) {
	var customer Customer

	query := `SELECT * FROM customers WHERE id = ?`

	err := r.db.GetContext(ctx, &customer, query, id)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]Customer, int, error) {
	var total int
	var customers []Customer

	baseWhere := ""
	var args []interface{}

	if strings.TrimSpace(params.Search) != "" {
		pattern := "%" + params.Search + "%"
		baseWhere = " WHERE (name LIKE ? OR phone LIKE ? OR email LIKE ? OR customer_code LIKE ?)"
		args = append(args, pattern, pattern, pattern, pattern)
	}

	countQuery := "SELECT COUNT(*) FROM customers" + baseWhere
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := "SELECT * FROM customers" + baseWhere + " ORDER BY id DESC LIMIT ? OFFSET ?"
	selectArgs := append(args, params.PerPage, params.Offset)

	err = r.db.SelectContext(ctx, &customers, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *Repository) Update(ctx context.Context, customer *Customer) error {
	query := `
		UPDATE customers
		SET name = ?, business_name = ?, email = ?, phone = ?, address = ?, city = ?, province = ?, notes = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		customer.Name,
		customer.BusinessName,
		customer.Email,
		customer.Phone,
		customer.Address,
		customer.City,
		customer.Province,
		customer.Notes,
		customer.ID,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM customers WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *Repository) GetNextCode(ctx context.Context) (string, error) {
	var maxCode *string

	query := `SELECT MAX(customer_code) FROM customers`
	err := r.db.GetContext(ctx, &maxCode, query)
	if err != nil {
		return "", err
	}

	nextNum := 1
	if maxCode != nil {
		parts := strings.Split(*maxCode, "-")
		if len(parts) == 2 {
			num, err := strconv.Atoi(parts[1])
			if err == nil {
				nextNum = num + 1
			}
		}
	}

	return fmt.Sprintf("CUST-%03d", nextNum), nil
}
