package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	qrcode "github.com/skip2/go-qrcode"

	"github.com/kittu533/reputap-be/internal/config"
	"github.com/kittu533/reputap-be/internal/database"
)

type BatchItem struct {
	Index        int
	SerialNumber string
	Code         string
	QRURL        string
	NFCURL       string
	ImagePath    string
}

func main() {
	countFlag := flag.Int("n", 20, "Number of blank products to produce")
	prefixFlag := flag.String("prefix", "RP26", "Serial number prefix")
	flag.Parse()

	count := *countFlag
	if count <= 0 {
		count = 20
	}

	cfg := config.Load()
	db, err := database.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	timestamp := time.Now().Format("20060102-150405")
	batchDir := filepath.Join("output", "production", fmt.Sprintf("batch_%d_%s", count, timestamp))
	qrDir := filepath.Join(batchDir, "qrcodes")

	if err := os.MkdirAll(qrDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Find the current highest serial number with this prefix to avoid conflicts
	var maxID int
	db.QueryRowContext(context.Background(), `
		SELECT COALESCE(COUNT(*), 0) FROM product_instances WHERE serial_number LIKE $1
	`, *prefixFlag+"-%").Scan(&maxID)

	fmt.Printf("🏭 Starting production batch: %d units (Prefix: %s, Start index: %d)\n", count, *prefixFlag, maxID+1)

	var items []BatchItem
	for i := 1; i <= count; i++ {
		num := maxID + i
		serial := fmt.Sprintf("%s-%04d", *prefixFlag, num)
		code := fmt.Sprintf("rpt-%04d", num)
		qrURL := fmt.Sprintf("%s/%s?src=qr", cfg.RedirectBaseURL, code)
		nfcURL := fmt.Sprintf("%s/%s?src=nfc", cfg.RedirectBaseURL, code)
		imgFileName := fmt.Sprintf("%s.png", serial)
		imgPath := filepath.Join(qrDir, imgFileName)

		// 1. Generate QR Code image (512x512 PNG, crisp for print)
		if err := qrcode.WriteFile(qrURL, qrcode.Medium, 512, imgPath); err != nil {
			log.Fatalf("Failed to generate QR code for %s: %v", serial, err)
		}

		// 2. Insert into database (Status: unassigned / ready-stock)
		var instanceID uint64
		err := db.QueryRowContext(context.Background(), `
			INSERT INTO product_instances (serial_number, nfc_status, production_status, qc_status)
			VALUES ($1, 'not_configured', 'completed', 'passed')
			ON CONFLICT (serial_number) DO UPDATE SET updated_at = CURRENT_TIMESTAMP
			RETURNING id
		`, serial).Scan(&instanceID)
		if err != nil {
			log.Fatalf("Failed to insert product instance %s: %v", serial, err)
		}

		_, err = db.ExecContext(context.Background(), `
			INSERT INTO qr_codes (product_instance_id, code, status)
			VALUES ($1, $2, 'unassigned')
			ON CONFLICT (code) DO UPDATE SET status = EXCLUDED.status
		`, instanceID, code)
		if err != nil {
			log.Fatalf("Failed to insert qr_code %s: %v", code, err)
		}

		items = append(items, BatchItem{
			Index:        i,
			SerialNumber: serial,
			Code:         code,
			QRURL:        qrURL,
			NFCURL:       nfcURL,
			ImagePath:    filepath.Join("qrcodes", imgFileName),
		})
	}

	// 3. Generate CSV Manifest for Digital Print Shop (Data Merge / VDP)
	csvPath := filepath.Join(batchDir, "manifest.csv")
	csvFile, err := os.Create(csvPath)
	if err != nil {
		log.Fatalf("Failed to create manifest.csv: %v", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write CSV Header
	writer.Write([]string{"No", "Serial_Number", "QR_Code", "QR_Redirect_URL", "QR_Image_Filename", "NFC_URL_Payload", "Status"})
	for _, it := range items {
		writer.Write([]string{
			strconv.Itoa(it.Index),
			it.SerialNumber,
			it.Code,
			it.QRURL,
			it.ImagePath,
			it.NFCURL,
			"READY_TO_ACTIVATE",
		})
	}

	// 4. Instructions file
	readmePath := filepath.Join(batchDir, "PANDUAN_PRODUKSI.txt")
	readmeContent := fmt.Sprintf(`=======================================================
REPUTAP PRODUCTION BATCH MANIFEST (%d UNITS)
Batch ID: batch_%d_%s
=======================================================

1. UNTUK PERCETAKAN (Stiker / Akrilik):
   - Gunakan file: manifest.csv
   - Kolom "QR_Image_Filename" berisi link gambar di folder "qrcodes/"
   - Di InDesign / Illustrator / CorelDraw: pakai fitur "Data Merge"
   - Cetak bahan: Vinyl Doff + Akrilik 2mm / 3mm

2. UNTUK OPERATOR NFC (Tulis ke Chip NTAG213):
   - Buka aplikasi "NFC Tools" (Android / iOS) atau pakai USB ACR122U
   - Menu: Write -> Add a record -> Custom URL / URI
   - Tulis URL yang tertera di kolom "NFC_URL_Payload" sesuai serial number
   - Tempel stiker NFC di bagian belakang akrilik

3. AKTIVASI SAAT TERJUAL:
   - Kartu ini statusnya "unassigned" (siap pakai).
   - Saat ada toko yang beli, tinggal buka dashboard:
     reputap.id/qr-codes -> Assign ke Toko & masukkan link Google Review.
   - Atau saat kartu di-tap sebelum aktif, akan otomatis muncul halaman
     aktivasi resmi Reputap!
`, count, count, timestamp)

	os.WriteFile(readmePath, []byte(readmeContent), 0644)

	fmt.Println("\n✅ 20 Units Produced Successfully!")
	fmt.Printf("📁 Output Folder : %s\n", batchDir)
	fmt.Printf("📊 CSV Manifest  : %s\n", csvPath)
	fmt.Printf("🖼️  QR Codes     : %s (%d high-res PNG files)\n", qrDir, count)
	fmt.Println("ℹ️  Semua unit tersimpan di database dengan status: 'unassigned' (siap diaktivasi)")
}
