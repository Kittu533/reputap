package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kittu533/reputap-be/internal/config"
	"github.com/kittu533/reputap-be/internal/database"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	fmt.Println("🧹 Resetting Reputap transactional data...")

	query := `
	BEGIN;

	TRUNCATE TABLE 
		scan_logs,
		audit_logs,
		quality_controls,
		production_jobs,
		design_assets,
		shipments,
		qr_codes,
		product_instances,
		order_status_history,
		order_items,
		orders,
		stores,
		customers
	RESTART IDENTITY CASCADE;

	DELETE FROM users WHERE email != 'admin@reputap.id';

	INSERT INTO users (name, email, password_hash, role_id, status)
	VALUES ('Super Admin', 'admin@reputap.id', '$2a$10$BWefn4Lz1nFfOA92nVSqFuOhUPCI/RXFxTEBz0l6G6hOxHPEBli7e', 1, 'active')
	ON CONFLICT (email) DO UPDATE SET status = 'active', role_id = 1;

	COMMIT;
	`

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("Failed to reset database: %v", err)
	}

	fmt.Println("✅ Database successfully formatted & reset to clean state!")
	fmt.Println("✨ Data transaksi, customer, store, order, QR, dan scan dibersihkan.")
	fmt.Println("👤 Admin default aktif:")
	fmt.Println("   Email    : admin@reputap.id")
	fmt.Println("   Password : admin123")
	fmt.Println("   Role     : super_admin")
}
