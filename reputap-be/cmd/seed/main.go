package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/kittu533/reputap-be/internal/config"
	"github.com/kittu533/reputap-be/internal/database"
)

func ctx() context.Context {
	return context.Background()
}

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	fmt.Println("🌱 Seeding database...")

	pwHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	hashStr := string(pwHash)

	// 1. Users
	var adminID, opID, mgrID uint64
	db.QueryRowContext(ctx(), `
		INSERT INTO users (name, email, password_hash, role_id, status)
		VALUES ('Super Admin', 'admin@reputap.id', $1, 1, 'active')
		ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, hashStr).Scan(&adminID)

	db.QueryRowContext(ctx(), `
		INSERT INTO users (name, email, password_hash, role_id, status)
		VALUES ('Siti Manager', 'manager@reputap.id', $1, 2, 'active')
		ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, hashStr).Scan(&mgrID)

	db.QueryRowContext(ctx(), `
		INSERT INTO users (name, email, password_hash, role_id, status)
		VALUES ('Budi Operator', 'operator@reputap.id', $1, 3, 'active')
		ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, hashStr).Scan(&opID)
	fmt.Println("  ✓ Users seeded")

	// 2. Customers
	type custData struct {
		code, name, bName, email, phone, addr, city, prov string
	}
	customers := []custData{
		{"CUST-001", "Budi Santoso", "Kopi Kenangan Senja", "budi@kopisenja.com", "081234567890", "Jl. Senopati No. 45", "Jakarta Selatan", "DKI Jakarta"},
		{"CUST-002", "Linda Wijaya", "Glow Aesthetic Clinic", "linda@glowclinic.id", "081987654321", "Jl. Raya Darmo No. 88", "Surabaya", "Jawa Timur"},
		{"CUST-003", "Hendra Pratama", "Barberhood Cuts", "hendra@barberhood.com", "082133445566", "Jl. Riau No. 12", "Bandung", "Jawa Barat"},
		{"CUST-004", "Maya Indah", "Rasa Nusantara Resto", "maya@rasanusantara.id", "085611223344", "Jl. Sunset Road No. 99", "Badung", "Bali"},
	}

	custIDs := make([]uint64, len(customers))
	for i, c := range customers {
		db.QueryRowContext(ctx(), `
			INSERT INTO customers (customer_code, name, business_name, email, phone, address, city, province)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (customer_code) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, c.code, c.name, c.bName, c.email, c.phone, c.addr, c.city, c.prov).Scan(&custIDs[i])
	}
	fmt.Println("  ✓ Customers seeded")

	// 3. Stores
	type storeData struct {
		custID uint64
		code, name, slug, cat, addr, gName, rURL string
		verified                                 bool
	}
	stores := []storeData{
		{custIDs[0], "STORE-001", "Kopi Senja Senopati", "kopi-senja-senopati", "coffee_shop", "Jl. Senopati No. 45, Jakarta", "Kopi Kenangan Senja Senopati", "https://g.page/r/kopi-senja/review", true},
		{custIDs[0], "STORE-002", "Kopi Senja Kemang", "kopi-senja-kemang", "coffee_shop", "Jl. Kemang Raya No. 10, Jakarta", "Kopi Kenangan Senja Kemang", "https://g.page/r/kopi-kemang/review", true},
		{custIDs[1], "STORE-003", "Glow Clinic Surabaya", "glow-clinic-sby", "clinic", "Jl. Raya Darmo No. 88, Surabaya", "Glow Aesthetic Clinic Darmo", "https://g.page/r/glow-clinic/review", true},
		{custIDs[2], "STORE-004", "Barberhood Riau", "barberhood-riau", "barbershop", "Jl. Riau No. 12, Bandung", "Barberhood Cuts & Grooming", "https://g.page/r/barberhood/review", false},
		{custIDs[3], "STORE-005", "Rasa Nusantara Seminyak", "rasa-nusantara-seminyak", "restaurant", "Jl. Sunset Road No. 99, Bali", "Rasa Nusantara Authentic Dining", "https://g.page/r/rasa-nusantara/review", true},
	}

	storeIDs := make([]uint64, len(stores))
	for i, s := range stores {
		db.QueryRowContext(ctx(), `
			INSERT INTO stores (customer_id, store_code, name, slug, category, address, google_business_name, review_url, review_url_verified, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'active')
			ON CONFLICT (store_code) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, s.custID, s.code, s.name, s.slug, s.cat, s.addr, s.gName, s.rURL, s.verified).Scan(&storeIDs[i])
	}
	fmt.Println("  ✓ Stores seeded")

	// 4. Products
	type prodData struct {
		sku, name, desc, mat string
		w, h, th             float64
		chip, finish         string
		cost, price          float64
	}
	prods := []prodData{
		{"PROD-AKRILIK-A6", "Stand Akrilik A6 NFC & QR", "Akrilik display A6 dengan built-in NFC chip dan cetak UV QR Code", "Akrilik Bening 3mm", 105.0, 148.0, 3.0, "NTAG213", "UV Print Matte", 25000, 75000},
		{"PROD-KAYU-MINI", "Wooden NFC Stand Mini", "Dudukan kayu jati elegan laser engraved dengan chip NFC", "Kayu Jati Solid", 80.0, 100.0, 15.0, "NTAG213", "Natural Teak Oil", 35000, 95000},
		{"PROD-CARD-PVC", "Smart Tap PVC Card", "Kartu NFC seukuran kartu ATM untuk ulasan cepat kasir & staff", "PVC Matte Finish", 85.6, 54.0, 0.8, "NTAG216", "Matte Lamination", 12000, 45000},
		{"PROD-METAL-PLATE", "Brushed Metal NFC Plate", "Plat logam stainless premium anti-cuaca tahan gores", "Stainless Steel 304", 100.0, 100.0, 1.5, "NTAG213", "Brushed Metal", 45000, 120000},
	}

	prodIDs := make([]uint64, len(prods))
	for i, p := range prods {
		db.QueryRowContext(ctx(), `
			INSERT INTO products (sku, name, description, material, width_mm, height_mm, thickness_mm, nfc_chip_type, finishing, production_cost, selling_price, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'active')
			ON CONFLICT (sku) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, p.sku, p.name, p.desc, p.mat, p.w, p.h, p.th, p.chip, p.finish, p.cost, p.price).Scan(&prodIDs[i])
	}
	fmt.Println("  ✓ Products seeded")

	// 5. Orders
	type ordData struct {
		num         string
		custID      uint64
		storeID     uint64
		status      string
		payStatus   string
		sub, dsc, t float64
	}
	orders := []ordData{
		{"ORD-2026-0001", custIDs[0], storeIDs[0], "completed", "paid", 225000, 0, 225000},
		{"ORD-2026-0002", custIDs[1], storeIDs[2], "production", "paid", 190000, 10000, 180000},
		{"ORD-2026-0003", custIDs[2], storeIDs[3], "qc", "paid", 95000, 0, 95000},
		{"ORD-2026-0004", custIDs[3], storeIDs[4], "pending", "unpaid", 240000, 0, 240000},
	}

	orderIDs := make([]uint64, len(orders))
	for i, o := range orders {
		db.QueryRowContext(ctx(), `
			INSERT INTO orders (order_number, customer_id, store_id, order_status, payment_status, subtotal, discount, total, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (order_number) DO UPDATE SET order_status = EXCLUDED.order_status
			RETURNING id
		`, o.num, o.custID, o.storeID, o.status, o.payStatus, o.sub, o.dsc, o.t, adminID).Scan(&orderIDs[i])
	}
	fmt.Println("  ✓ Orders seeded")

	// 6. Order Items
	type itemData struct {
		ordID, prodID uint64
		qty           int
		price, sub    float64
	}
	items := []itemData{
		{orderIDs[0], prodIDs[0], 3, 75000, 225000},
		{orderIDs[1], prodIDs[1], 2, 95000, 190000},
		{orderIDs[2], prodIDs[1], 1, 95000, 95000},
		{orderIDs[3], prodIDs[3], 2, 120000, 240000},
	}

	itemIDs := make([]uint64, len(items))
	for i, it := range items {
		// Clean existing items for the order if re-running
		var existingID uint64
		err := db.QueryRowContext(ctx(), `SELECT id FROM order_items WHERE order_id = $1 AND product_id = $2`, it.ordID, it.prodID).Scan(&existingID)
		if err == sql.ErrNoRows {
			db.QueryRowContext(ctx(), `
				INSERT INTO order_items (order_id, product_id, quantity, unit_price, subtotal)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id
			`, it.ordID, it.prodID, it.qty, it.price, it.sub).Scan(&itemIDs[i])
		} else {
			itemIDs[i] = existingID
		}
	}
	fmt.Println("  ✓ Order items seeded")

	// 7. Product Instances & QR Codes
	type instData struct {
		itemID, storeID uint64
		serial, qrCode  string
		nfc, prod, qc   string
	}
	instances := []instData{
		{itemIDs[0], storeIDs[0], "SN-2026-0001", "kps-seno-1", "configured", "completed", "passed"},
		{itemIDs[0], storeIDs[0], "SN-2026-0002", "kps-seno-2", "configured", "completed", "passed"},
		{itemIDs[0], storeIDs[0], "SN-2026-0003", "kps-seno-3", "configured", "completed", "passed"},
		{itemIDs[1], storeIDs[2], "SN-2026-0004", "glw-sby-1", "configured", "in_progress", "pending"},
		{itemIDs[2], storeIDs[3], "SN-2026-0005", "brb-riau-1", "not_configured", "completed", "pending"},
	}

	qrIDs := make([]uint64, len(instances))
	for i, inst := range instances {
		var instID uint64
		err := db.QueryRowContext(ctx(), `
			INSERT INTO product_instances (order_item_id, store_id, serial_number, nfc_status, production_status, qc_status)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (serial_number) DO UPDATE SET nfc_status = EXCLUDED.nfc_status
			RETURNING id
		`, inst.itemID, inst.storeID, inst.serial, inst.nfc, inst.prod, inst.qc).Scan(&instID)
		if err != nil {
			continue
		}

		var qID uint64
		db.QueryRowContext(ctx(), `
			INSERT INTO qr_codes (product_instance_id, store_id, code, status)
			VALUES ($1, $2, $3, 'active')
			ON CONFLICT (code) DO UPDATE SET status = EXCLUDED.status
			RETURNING id
		`, instID, inst.storeID, inst.qrCode).Scan(&qID)
		qrIDs[i] = qID
	}
	fmt.Println("  ✓ Product instances & QR codes seeded")

	// 8. Scan Logs (for analytics graph)
	var scanCount int
	db.QueryRowContext(ctx(), `SELECT COUNT(*) FROM scan_logs`).Scan(&scanCount)
	if scanCount < 50 {
		sources := []string{"qr", "nfc"}
		devices := []string{"android", "ios"}
		userAgents := []string{
			"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Mobile/15E148 Safari/604.1",
			"Mozilla/5.0 (Linux; Android 14; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.119 Mobile Safari/537.36",
		}

		now := time.Now()
		// Generate scans across the last 30 days
		for day := 30; day >= 0; day-- {
			numScans := rand.Intn(12) + 2 // 2 to 14 scans per day
			targetDate := now.AddDate(0, 0, -day)

			for s := 0; s < numScans; s++ {
				qrID := qrIDs[rand.Intn(len(qrIDs))]
				if qrID == 0 {
					continue
				}
				source := sources[rand.Intn(len(sources))]
				device := devices[rand.Intn(len(devices))]
				ua := userAgents[rand.Intn(len(userAgents))]
				hour := rand.Intn(14) + 8 // Between 8am and 10pm
				min := rand.Intn(60)
				scanTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), hour, min, rand.Intn(60), 0, targetDate.Location())

				db.ExecContext(ctx(), `
					INSERT INTO scan_logs (qr_code_id, source, device_type, user_agent, referrer, scanned_at)
					VALUES ($1, $2, $3, $4, 'https://reputap.id/r', $5)
				`, qrID, source, device, ua, scanTime)
			}
		}
		fmt.Println("  ✓ Scan logs seeded (realistic 30-day timeline)")
	} else {
		fmt.Println("  ✓ Scan logs already populated")
	}

	// 9. Production Jobs & Quality Controls
	var instance1ID uint64
	db.QueryRowContext(ctx(), `SELECT id FROM product_instances WHERE serial_number = 'SN-2026-0001'`).Scan(&instance1ID)
	if instance1ID > 0 {
		var jobID uint64
		db.QueryRowContext(ctx(), `SELECT id FROM production_jobs WHERE product_instance_id = $1`, instance1ID).Scan(&jobID)
		if jobID == 0 {
			db.ExecContext(ctx(), `
				INSERT INTO production_jobs (product_instance_id, assigned_to, status, checklist_data, started_at, completed_at, notes)
				VALUES ($1, $2, 'completed', '{"print_cut": true, "nfc_embedded": true, "surface_clean": true}', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day', 'Finished UV printing & acrylic cutting')
			`, instance1ID, opID)

			db.ExecContext(ctx(), `
				INSERT INTO quality_controls (product_instance_id, inspected_by, result, checklist_data, notes, inspected_at)
				VALUES ($1, $2, 'passed', '{"nfc_tested": true, "url_redirect_verified": true, "aesthetic_check": true}', 'NFC reads instantly within 3cm', NOW() - INTERVAL '1 day')
			`, instance1ID, mgrID)
		}
	}
	fmt.Println("  ✓ Production jobs & QC records seeded")

	// 10. Shipment
	var shipID uint64
	db.QueryRowContext(ctx(), `SELECT id FROM shipments WHERE order_id = $1`, orderIDs[0]).Scan(&shipID)
	if shipID == 0 && orderIDs[0] > 0 {
		db.ExecContext(ctx(), `
			INSERT INTO shipments (order_id, recipient_name, recipient_address, recipient_phone, courier, tracking_number, shipping_status, shipped_at, delivered_at)
			VALUES ($1, 'Budi Santoso', 'Jl. Senopati No. 45, Jakarta Selatan', '081234567890', 'JNE Express', 'JNE889102931', 'delivered', NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 day')
		`, orderIDs[0])
	}
	fmt.Println("  ✓ Shipments seeded")

	// 11. Audit Logs
	var auditCount int
	db.QueryRowContext(ctx(), `SELECT COUNT(*) FROM audit_logs`).Scan(&auditCount)
	if auditCount == 0 {
		db.ExecContext(ctx(), `
			INSERT INTO audit_logs (user_id, action, entity_type, entity_id, old_values, new_values, created_at) VALUES
			($1, 'USER_LOGIN', 'users', $1, NULL, '{"ip": "127.0.0.1"}', NOW() - INTERVAL '4 hours'),
			($1, 'ORDER_CREATED', 'orders', $2, NULL, '{"order_number": "ORD-2026-0001", "total": 225000}', NOW() - INTERVAL '3 days'),
			($3, 'PRODUCTION_COMPLETED', 'production_jobs', 1, '{"status": "in_progress"}', '{"status": "completed"}', NOW() - INTERVAL '1 day'),
			($4, 'QC_PASSED', 'quality_controls', 1, NULL, '{"result": "passed"}', NOW() - INTERVAL '1 day')
		`, adminID, orderIDs[0], opID, mgrID)
		fmt.Println("  ✓ Audit logs seeded")
	}

	fmt.Println("✨ Database seeding complete successfully!")
}
