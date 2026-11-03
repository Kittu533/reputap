# REPUTAP
## Product Requirements Document (PRD)

**Product Version:** 2.0  
**Document Status:** Product Definition & Development Blueprint  
**Product Category:** Digital Business Tools / Review Management  
**Business Model:** Done-for-You Product & Service  
**Frontend:** Nuxt 4  
**Backend:** Go (Golang)  
**Database:** MySQL 8  
**UI Framework:** Tailwind CSS + Nuxt UI  
**Architecture:** Modular Monolith  
**Initial Deployment:** Local Development  
**Target Market:** UMKM, coffee shops, restaurants, salons, barbershops, retail stores, and other local businesses.

---

# 1. EXECUTIVE SUMMARY

## 1.1 Product Overview

REPUTAP adalah produk dan layanan yang menggabungkan perangkat fisik berupa display akrilik dengan teknologi NFC dan Dynamic QR Code untuk mempermudah pelanggan mengakses halaman Google Review milik suatu bisnis.

REPUTAP menggunakan model bisnis Done-for-You (DFY), yaitu seluruh proses konfigurasi, personalisasi, produksi, pengujian, dan persiapan produk dilakukan oleh tim REPUTAP.

Pemilik bisnis tidak perlu melakukan konfigurasi teknis, membuat QR Code, mengatur NFC, maupun menggunakan dashboard aplikasi.

Pelanggan cukup melakukan pemesanan dan memberikan informasi bisnis yang diperlukan. Tim REPUTAP akan menyiapkan produk sampai siap digunakan.

## 1.2 Core Concept

REPUTAP mengintegrasikan tiga komponen utama:

1. **Physical Product:** Display akrilik custom sebagai media interaksi pelanggan.
2. **Dynamic QR System:** QR Code unik yang dapat mempertahankan alamat permanen meskipun URL tujuan berubah.
3. **NFC Technology:** Teknologi komunikasi jarak dekat yang memungkinkan pelanggan membuka halaman review hanya dengan mendekatkan smartphone.

## 1.3 Product Philosophy

Prinsip utama REPUTAP:

**Buy. Receive. Display. Get Reviewed.**

Pemilik bisnis hanya perlu:

- Membeli produk.
- Memberikan informasi bisnis.
- Menerima produk yang sudah dikonfigurasi.
- Meletakkan produk di lokasi strategis.
- Memantau ulasan yang diberikan pelanggan melalui Google Maps/Google Business Profile.

Seluruh proses teknis ditangani oleh tim REPUTAP.

## 1.4 Problem Statement

Beberapa permasalahan yang ingin diselesaikan:

- Pelanggan kesulitan menemukan halaman Google Review bisnis.
- Pemilik bisnis harus menjelaskan cara memberikan ulasan secara manual.
- QR Code statis tidak fleksibel ketika URL tujuan berubah.
- Pengelolaan QR dan NFC membutuhkan pengetahuan teknis.
- Produk review yang tersedia belum tentu menyediakan layanan konfigurasi dan pengelolaan yang praktis.
- Pemilik bisnis membutuhkan media review yang menarik dan sesuai identitas brand mereka.

## 1.5 Proposed Solution

REPUTAP menyediakan display akrilik NFC + QR yang sudah dikonfigurasi secara khusus untuk setiap bisnis.

Setiap unit memiliki identitas digital permanen yang terhubung dengan sistem Dynamic QR REPUTAP.

Tim dapat memperbarui URL Google Review melalui internal dashboard tanpa harus mencetak ulang QR Code atau mengubah konfigurasi NFC.

---

# 2. PRODUCT VISION & OBJECTIVES

## 2.1 Product Vision

Menjadi solusi praktis bagi UMKM untuk mempermudah pelanggan mengakses halaman ulasan bisnis melalui integrasi teknologi fisik dan digital.

## 2.2 Product Mission

1. Menyediakan media pengumpulan ulasan yang mudah digunakan.
2. Menghilangkan kerumitan konfigurasi teknologi bagi pemilik bisnis.
3. Mengotomatisasi proses pengelolaan pesanan dan konfigurasi produk.
4. Menyediakan sistem QR dan NFC yang fleksibel.
5. Mengembangkan layanan fisik menjadi ekosistem produk digital.

## 2.3 Business Objectives

| ID | Objective | Target Awal |
|---|---|---|
| BO-01 | Memvalidasi kebutuhan pasar | 5–10 UMKM |
| BO-02 | Memvalidasi produk fisik | 2 prototype |
| BO-03 | Menguji sistem Dynamic QR | Berhasil pada seluruh unit uji |
| BO-04 | Meminimalkan kesalahan produksi | 0 kesalahan konfigurasi saat QC |
| BO-05 | Memvalidasi penjualan | 5 transaksi berbayar |
| BO-06 | Memastikan operasional terstruktur | Seluruh order memiliki status |

Target tersebut merupakan sasaran market validation, bukan jaminan pencapaian.

## 2.4 Product Success Metrics

- Jumlah pesanan yang masuk.
- Jumlah pesanan yang berhasil diselesaikan.
- Rata-rata waktu penyelesaian pesanan.
- Persentase kesalahan konfigurasi.
- Jumlah QR aktif.
- Total QR scan.
- Total NFC tap.
- Redirect success rate.
- Jumlah komplain produk.
- Persentase pelanggan yang melakukan repeat order.

---

# 3. BUSINESS MODEL

## 3.1 Business Model Definition

REPUTAP menggunakan pendekatan Done-for-You Product.

Tim REPUTAP bertindak sebagai penyedia produk sekaligus operator sistem.

Pelanggan tidak membeli akses dashboard. Pelanggan membeli produk fisik yang sudah siap digunakan.

## 3.2 Customer Journey

1. Pelanggan menemukan produk REPUTAP.
2. Pelanggan melakukan pemesanan.
3. Pelanggan memberikan informasi bisnis.
4. Tim REPUTAP melakukan verifikasi data.
5. Tim mengonfigurasi QR dan NFC.
6. Tim membuat desain custom.
7. Tim melakukan produksi fisik.
8. Tim menjalankan Quality Control.
9. Produk dikirim kepada pelanggan.
10. Pelanggan meletakkan display di toko.
11. Pengunjung melakukan scan atau tap.
12. Pengunjung diarahkan ke Google Review bisnis.

## 3.3 Revenue Model

Sumber pendapatan awal:

- Penjualan display akrilik NFC + QR.
- Penjualan paket tambahan display.
- Custom branding.
- Penggantian atau penambahan produk.

Potensi pendapatan berikutnya:

- Paket pemeliharaan.
- Pengelolaan beberapa cabang.
- Advanced analytics.
- Paket bisnis dengan beberapa display.

Fitur berlangganan bukan bagian dari MVP.

## 3.4 Target Customers

Primary market:
- Coffee shop.
- Restoran.
- Barbershop.
- Salon kecantikan.
- Klinik.
- Toko retail.
- Hotel dan penginapan.
- Tempat usaha berbasis layanan.

Secondary market:
- Agen pemasaran lokal.
- Konsultan bisnis.
- Penyedia jasa branding UMKM.

---

# 4. PRODUCT SCOPE

## 4.1 Scope Overview

Produk dibagi menjadi tiga bagian utama.

### A. Physical Product

Display akrilik yang dilengkapi QR Code dan NFC.

### B. Digital System

Sistem internal yang mengelola data pelanggan, pesanan, QR, NFC, URL tujuan, dan statistik penggunaan.

### C. Operational Management

Workflow internal yang membantu tim menangani pesanan sampai produk selesai diproduksi dan dikirim.

## 4.2 Physical Product Specification

Spesifikasi prototype pertama:

| Komponen | Spesifikasi |
|---|---|
| Product Name | REPUTAP Smart Review Stand |
| Panel | Akrilik transparan |
| Material Thickness | 2 mm |
| Panel Width | 10 cm |
| Panel Height | 15 cm |
| Main Structure | Satu lembar akrilik |
| Stand Type | L-shape / kaki penyangga terpisah |
| Visual | Custom vinyl sticker |
| Finishing | Matte / doff |
| NFC Chip | NTAG213 |
| QR Technology | Dynamic QR |
| Product Orientation | Portrait |

Spesifikasi konstruksi final harus mengikuti pengujian stabilitas, metode bending atau base support, dan posisi komponen NFC.

## 4.3 Digital Product Specification

Sistem harus menyediakan:

- Admin authentication.
- Dashboard overview.
- Order management.
- Customer management.
- Business/store management.
- QR code management.
- NFC configuration management.
- Dynamic redirect engine.
- Design asset management.
- Production tracking.
- Quality control.
- Shipping management.
- Scan analytics.
- Activity logs.

---

# 5. USER ROLES & ACCESS CONTROL

## 5.1 Administrator

Pengguna internal yang mengelola seluruh operasional.

Permissions:
- Mengelola toko.
- Mengelola pesanan.
- Mengatur QR dan NFC.
- Mengubah URL Google Review.
- Mengelola desain.
- Memantau produksi.
- Melakukan QC.
- Mengelola pengiriman.
- Melihat statistik.
- Mengelola konfigurasi sistem.

## 5.2 Operator

Pengguna internal yang menangani proses operasional.

Permissions:
- Melihat pesanan yang ditugaskan.
- Memperbarui status produksi.
- Mengunggah desain.
- Menginput konfigurasi.
- Melakukan QC.
- Mengisi informasi pengiriman.

Operator tidak boleh mengubah konfigurasi keamanan dan hak akses administrator.

## 5.3 Super Admin

Pemilik sistem dengan akses penuh.

Permissions:
- Semua akses administrator.
- Manajemen pengguna.
- Pengaturan role dan permission.
- Pengaturan sistem.
- Pengelolaan audit log.
- Pengaturan domain dan redirect.

## 5.4 Customer

Customer tidak membutuhkan akun untuk MVP.

Customer hanya perlu:
- Melakukan pemesanan.
- Memberikan informasi bisnis.
- Menerima produk.
- Menggunakan produk.

---

# 6. FUNCTIONAL REQUIREMENTS

# MODULE 01 — AUTHENTICATION & AUTHORIZATION

## FR-AUTH-001: Admin Login

Sistem menyediakan halaman login internal.

Input:
- Email.
- Password.

Acceptance Criteria:
- Email harus memiliki format valid.
- Password diverifikasi menggunakan hash yang aman.
- Kredensial salah menghasilkan pesan kesalahan.
- Admin yang tidak terautentikasi tidak dapat membuka dashboard.
- Login berhasil membuat sesi autentikasi.
- Sistem mencatat waktu login.

## FR-AUTH-002: Logout

Admin dapat mengakhiri sesi.

Acceptance Criteria:
- Token atau sesi menjadi tidak valid.
- Pengguna diarahkan ke halaman login.
- Endpoint terlindungi menolak sesi yang telah berakhir.

## FR-AUTH-003: Role Management

Super Admin dapat mengatur role dan permission.

## FR-AUTH-004: Password Management

Admin dapat mengganti password.

Sistem harus menerapkan validasi password dan proteksi terhadap percobaan login berulang.

---

# MODULE 02 — DASHBOARD OVERVIEW

## FR-DASH-001: Summary Cards

Dashboard menampilkan:

- Total orders.
- Pending orders.
- Orders in production.
- Completed orders.
- Total customers.
- Total stores.
- Active QR codes.
- Total scans.
- Redirect success rate.

## FR-DASH-002: Order Analytics

Dashboard menampilkan grafik pesanan berdasarkan:

- Hari.
- Minggu.
- Bulan.
- Status pesanan.

## FR-DASH-003: Scan Analytics

Dashboard menampilkan:
- Total QR scans.
- Total NFC taps.
- Unique visitor estimate.
- Scan per hari.
- Scan per toko.
- Scan per produk.

## FR-DASH-004: Recent Activities

Menampilkan aktivitas terbaru:
- Pesanan baru.
- Perubahan URL.
- QR baru dibuat.
- NFC berhasil dikonfigurasi.
- Produksi dimulai.
- QC selesai.
- Produk dikirim.

---

# MODULE 03 — CUSTOMER MANAGEMENT

Modul untuk mengelola informasi pemilik bisnis yang membeli REPUTAP.

## FR-CUST-001: Create Customer

Admin dapat membuat data customer.

Fields:
- Customer ID.
- Nama pemilik.
- Nama bisnis.
- Email.
- Nomor WhatsApp.
- Alamat.
- Kota.
- Provinsi.
- Catatan internal.

## FR-CUST-002: Update Customer

Admin dapat mengubah data pelanggan.

## FR-CUST-003: Customer Detail

Halaman detail pelanggan menampilkan:
- Profil pelanggan.
- Daftar bisnis.
- Riwayat pesanan.
- Produk yang dimiliki.
- Riwayat komunikasi.
- Catatan operasional.

## FR-CUST-004: Customer Search

Admin dapat mencari pelanggan menggunakan:
- Nama.
- Nomor WhatsApp.
- Email.
- Customer ID.

---

# MODULE 04 — BUSINESS / STORE MANAGEMENT

Modul untuk mengelola data bisnis pelanggan.

## FR-STORE-001: Create Store

Admin dapat mendaftarkan bisnis baru.

Fields:
- Store ID.
- Store name.
- Business category.
- Owner.
- Google Business Profile name.
- Google Review URL.
- Address.
- Contact information.
- Status.

## FR-STORE-002: Business Verification

Admin dapat memverifikasi bahwa URL Google Review benar-benar mengarah ke bisnis yang dimaksud.

Acceptance Criteria:
- URL tidak boleh kosong.
- URL harus menggunakan HTTPS.
- URL harus mengarah ke domain Google yang diperbolehkan.
- Admin dapat membuka preview tujuan.
- Admin menandai verifikasi berhasil atau gagal.

## FR-STORE-003: Multi Location

Satu customer dapat memiliki beberapa cabang bisnis.

Contoh:

Customer: PT Kopi Nusantara

- Kopi Senja Cirebon.
- Kopi Senja Bandung.
- Kopi Senja Jakarta.

Setiap cabang dapat memiliki URL Google Review berbeda.

## FR-STORE-004: Store Status

Status:
- Active.
- Inactive.
- Suspended.

Store yang nonaktif tidak dapat menerima redirect.

---

# MODULE 05 — ORDER MANAGEMENT

Modul untuk menangani seluruh proses pemesanan.

## FR-ORDER-001: Create Order

Admin dapat membuat pesanan baru.

Fields:
- Order ID.
- Customer.
- Business.
- Product.
- Quantity.
- Unit price.
- Discount.
- Total amount.
- Payment status.
- Order status.
- Notes.

## FR-ORDER-002: Order Detail

Halaman detail pesanan menampilkan:
- Informasi customer.
- Informasi toko.
- Produk yang dibeli.
- Link Google Review.
- Konfigurasi QR.
- Konfigurasi NFC.
- Desain produk.
- Status produksi.
- QC results.
- Informasi pengiriman.
- Riwayat perubahan.

## FR-ORDER-003: Order Status

Workflow:

1. NEW_ORDER
2. WAITING_CUSTOMER_DATA
3. DATA_VERIFICATION
4. CONFIGURATION
5. DESIGN
6. DESIGN_APPROVAL
7. PRODUCTION
8. QUALITY_CONTROL
9. READY_TO_SHIP
10. SHIPPED
11. COMPLETED

Status tambahan:
- CANCELLED.
- ON_HOLD.
- REVISION_REQUIRED.

Setiap perubahan status harus mencatat pengguna, waktu, dan catatan.

## FR-ORDER-004: Order Validation

Pesanan tidak dapat masuk produksi jika:
- Data bisnis belum lengkap.
- URL Google Review belum diverifikasi.
- Desain belum disetujui jika approval diperlukan.
- Pembayaran belum memenuhi kebijakan internal.

## FR-ORDER-005: Order Search & Filter

Filter berdasarkan:
- Order ID.
- Customer.
- Store.
- Status.
- Payment status.
- Date range.

## FR-ORDER-006: Order History

Sistem menyimpan seluruh perubahan pesanan.

---

# MODULE 06 — PRODUCT CATALOG

Modul untuk mengelola produk REPUTAP.

## FR-PROD-001: Product Management

Admin dapat menambahkan produk.

Contoh:
- REPUTAP Acrylic Stand.
- REPUTAP Acrylic Stand Premium.
- REPUTAP Multi Display Pack.

## FR-PROD-002: Product Specification

Fields:
- Product name.
- SKU.
- Description.
- Material.
- Dimensions.
- Thickness.
- NFC chip type.
- Finishing.
- Production cost.
- Selling price.
- Product status.

## FR-PROD-003: Product Pricing

Sistem menyimpan:
- Harga jual.
- Biaya bahan.
- Biaya cetak.
- Biaya perakitan.
- Packaging.
- Biaya tambahan.
- Estimasi margin.

---

# MODULE 07 — DYNAMIC QR MANAGEMENT

Modul inti REPUTAP.

## FR-QR-001: Generate QR Code

Setiap produk memiliki identifier unik.

Contoh:
- RP-000001.
- RP-000002.
- RP-000003.

URL permanen:

https://reputap.id/r/RP-000001

QR Code harus dibuat dari URL permanen tersebut.

## FR-QR-002: QR Preview

Admin dapat melihat preview QR.

Informasi:
- QR image.
- QR ID.
- Product ID.
- Store.
- Destination URL.
- Creation date.
- Status.

## FR-QR-003: QR Download

QR dapat diunduh dalam:
- PNG.
- SVG.

## FR-QR-004: QR Reassignment

Satu QR tidak boleh dipindahkan ke bisnis lain tanpa proses reassignment yang tercatat.

Jika produk diganti atau pelanggan berganti kepemilikan, admin harus memproses reassignment secara eksplisit.

## FR-QR-005: QR Status

Status:
- Active.
- Inactive.
- Suspended.
- Retired.

QR yang inactive tidak boleh diarahkan ke URL pelanggan.

## FR-QR-006: Permanent QR

URL QR harus tetap sama selama masa penggunaan produk.

Perubahan URL Google Review tidak mengubah QR fisik.

## FR-QR-007: QR Source Tracking

QR menggunakan URL dengan source identifier.

Contoh:

https://reputap.id/r/RP-000001?src=qr

Dengan demikian sistem dapat membedakan scan QR dan NFC.

---

# MODULE 08 — NFC CONFIGURATION

## FR-NFC-001: NFC Configuration

Sistem menghasilkan URL permanen yang akan ditulis ke NFC tag.

Contoh:

https://reputap.id/r/RP-000001?src=nfc

## FR-NFC-002: NFC Encoding

Operator menggunakan smartphone atau perangkat NFC writer untuk menulis URL ke NTAG213.

Sistem menyediakan:
- NFC URL.
- Product ID.
- QR ID.
- Encoding status.
- Encoding date.
- Operator.

## FR-NFC-003: NFC Testing

Setelah proses encoding, operator wajib melakukan pengujian.

Acceptance Criteria:
- Tag dapat dibaca.
- URL sesuai dengan produk.
- URL dapat dibuka.
- Redirect menuju Google Review yang benar.

## FR-NFC-004: NFC Status

Status:
- Not Configured.
- Configured.
- Tested.
- Failed.
- Replaced.

## FR-NFC-005: NFC Re-encoding

Jika URL permanen REPUTAP tidak berubah, NFC tidak perlu ditulis ulang ketika Google Review URL berubah.

---

# MODULE 09 — REDIRECT ENGINE

Modul backend untuk menangani QR dan NFC.

## FR-REDIRECT-001: Resolve Identifier

Endpoint:

GET /r/{code}

Sistem mencari QR Code berdasarkan identifier.

## FR-REDIRECT-002: Resolve Destination

Sistem mengambil URL tujuan yang aktif berdasarkan relasi QR dan produk.

## FR-REDIRECT-003: Redirect

Jika data valid, backend memberikan HTTP 302 redirect menuju Google Review.

## FR-REDIRECT-004: Invalid Identifier

Jika identifier tidak ditemukan, sistem menampilkan halaman:
- QR tidak ditemukan.
- QR sudah tidak aktif.
- Hubungi tim REPUTAP.

## FR-REDIRECT-005: Destination Validation

URL tujuan harus melewati validasi.

Hanya domain Google yang diizinkan dalam konfigurasi awal.

Sistem tidak boleh menjadi open redirect bebas.

## FR-REDIRECT-006: Redirect Logging

Setiap request dicatat dengan informasi:
- QR ID.
- Source.
- Timestamp.
- User agent.
- Device category.
- Referrer jika tersedia.
- Anonymous identifier jika mekanisme privasi mengizinkan.

IP address tidak perlu disimpan secara permanen dalam MVP.

## FR-REDIRECT-007: Redirect Performance

Redirect harus memiliki proses yang ringan.

Target internal:
- P95 response time maksimal 500 ms dalam lingkungan production normal.
- Tidak memerlukan rendering frontend Nuxt.
- Tidak memerlukan autentikasi pengunjung.

---

# MODULE 10 — DESIGN & ASSET MANAGEMENT

Modul untuk mengelola desain visual produk.

## FR-DESIGN-001: Design Template

Admin dapat menyimpan template desain.

Kategori:
- Coffee Shop.
- Restaurant.
- Barbershop.
- Beauty Salon.
- Clinic.
- Retail.
- General Business.

## FR-DESIGN-002: Custom Design

Setiap pesanan dapat memiliki desain khusus.

Fields:
- Design ID.
- Order ID.
- Store ID.
- Template.
- Logo.
- Brand color.
- Business name.
- QR asset.
- Preview image.
- Revision number.
- Approval status.

## FR-DESIGN-003: Asset Upload

Admin dapat mengunggah:
- PNG.
- JPG.
- SVG.
- PDF.

Sistem membatasi ukuran file dan memvalidasi tipe file.

## FR-DESIGN-004: Design Revision

Sistem menyimpan versi desain.

Contoh:
- Version 1.
- Version 2.
- Final Version.

## FR-DESIGN-005: Design Approval

Status:
- Draft.
- Pending Approval.
- Approved.
- Revision Required.

Desain yang memerlukan persetujuan pelanggan harus disetujui sebelum masuk produksi.

## FR-DESIGN-006: Print Export

Sistem menyimpan file desain final yang siap digunakan untuk produksi.

Spesifikasi:
- Ukuran sesuai panel fisik.
- Resolusi cetak memadai.
- QR memiliki quiet zone.
- Kontras QR memadai.
- QR dapat dipindai pada ukuran cetak final.

---

# MODULE 11 — PRODUCTION MANAGEMENT

## FR-PROD-OPS-001: Production Queue

Menampilkan pesanan yang siap diproduksi.

## FR-PROD-OPS-002: Production Checklist

Checklist:
- Akrilik tersedia.
- Cutting selesai.
- Finishing selesai.
- Vinyl selesai dicetak.
- NFC tersedia.
- QR tersedia.
- Assembly selesai.
- Packaging selesai.

## FR-PROD-OPS-003: Production Assignment

Admin dapat menentukan operator produksi.

## FR-PROD-OPS-004: Production Notes

Operator dapat menambahkan catatan.

Contoh:
- Ukuran cutting tidak sesuai.
- Stiker perlu dicetak ulang.
- NFC sticker gagal terbaca.

## FR-PROD-OPS-005: Production Cost

Sistem menyimpan biaya aktual per pesanan.

---

# MODULE 12 — QUALITY CONTROL

## FR-QC-001: QC Checklist

Setiap produk harus melewati pengujian.

Checklist:

1. Ukuran akrilik sesuai spesifikasi.
2. Tidak ada retak atau kerusakan.
3. Finishing rapi.
4. Stiker vinyl tidak menggelembung.
5. QR Code dapat dipindai.
6. QR mengarah ke toko yang benar.
7. NFC dapat dibaca.
8. NFC mengarah ke toko yang benar.
9. QR dan NFC menggunakan product identifier yang sama.
10. Produk stabil saat diletakkan.
11. Packaging aman.

## FR-QC-002: QC Result

Status:
- Passed.
- Failed.
- Rework Required.

Produk yang gagal QC tidak boleh masuk ke status Ready to Ship.

## FR-QC-003: QC Evidence

Operator dapat mengunggah foto produk dan hasil pengujian.

## FR-QC-004: QC History

Sistem menyimpan riwayat pemeriksaan.

---

# MODULE 13 — SHIPPING MANAGEMENT

## FR-SHIP-001: Shipping Information

Fields:
- Recipient.
- Address.
- Courier.
- Tracking number.
- Shipping date.
- Shipping status.

## FR-SHIP-002: Shipping Status

Status:
- Waiting Shipment.
- Packed.
- Shipped.
- Delivered.
- Delivery Failed.

## FR-SHIP-003: Shipment Tracking

Admin dapat menyimpan nomor resi.

Integrasi kurir belum wajib pada MVP.

---

# MODULE 14 — SCAN ANALYTICS

## FR-ANALYTICS-001: Scan Count

Menghitung total request QR dan NFC.

## FR-ANALYTICS-002: Source Analytics

Membedakan:
- QR Scan.
- NFC Tap.

Pembedaan dilakukan melalui parameter source pada URL.

## FR-ANALYTICS-003: Device Analytics

Mengelompokkan:
- Android.
- iOS.
- Desktop.
- Other.

## FR-ANALYTICS-004: Time Analytics

Statistik berdasarkan:
- Daily.
- Weekly.
- Monthly.

## FR-ANALYTICS-005: Store Analytics

Admin dapat membandingkan penggunaan antar bisnis.

## FR-ANALYTICS-006: Privacy

Statistik harus memperhatikan privasi pengguna.

Jangan menyimpan data personal pengunjung yang tidak diperlukan.

Unique visitor hanya merupakan estimasi, bukan identitas orang yang terverifikasi.

---

# MODULE 15 — ACTIVITY LOG & AUDIT TRAIL

## FR-AUDIT-001: Activity Logging

Sistem mencatat:
- Login.
- Create order.
- Update order.
- Change review URL.
- Generate QR.
- NFC configuration.
- Production update.
- QC result.
- Shipping update.

## FR-AUDIT-002: Change History

Perubahan penting harus memiliki:
- User ID.
- Timestamp.
- Previous value.
- New value.
- Reason atau notes.

## FR-AUDIT-003: Audit Protection

Log aktivitas tidak boleh diubah oleh operator biasa.

---

# 7. BUSINESS PROCESS WORKFLOW

## 7.1 Order-to-Delivery Workflow

### Stage 1 — Order Creation

Pesanan masuk melalui channel penjualan.

Channel awal:
- WhatsApp.
- Instagram.
- Penjualan langsung.

Admin memasukkan pesanan ke sistem.

Output:
- Order ID.
- Customer record.
- Store record.
- Product item.

### Stage 2 — Customer Data Collection

Tim meminta:
- Nama bisnis.
- Logo.
- Alamat.
- Link Google Review.
- Identitas visual jika tersedia.
- Preferensi desain.

Order berpindah ke WAITING_CUSTOMER_DATA sampai seluruh data tersedia.

### Stage 3 — Data Verification

Tim memastikan link benar.

Verifikasi meliputi:
- URL dapat dibuka.
- URL menuju profil bisnis yang tepat.
- Profil bisnis sesuai nama pelanggan.
- Tidak ada kesalahan tujuan.

### Stage 4 — QR & NFC Configuration

Sistem menghasilkan identifier permanen.

Contoh:
RP-000001.

QR URL:
https://reputap.id/r/RP-000001?src=qr

NFC URL:
https://reputap.id/r/RP-000001?src=nfc

### Stage 5 — Design

Tim membuat banner custom berdasarkan template REPUTAP.

### Stage 6 — Production

Tim melakukan:
- Cutting akrilik.
- Cetak stiker.
- Pemasangan NFC.
- Pemasangan desain.
- Finishing.
- Assembly.

### Stage 7 — Quality Control

Semua fungsi diuji.

### Stage 8 — Shipping

Produk dikemas dan dikirim.

### Stage 9 — Completion

Order ditandai selesai setelah pengiriman atau penerimaan dikonfirmasi sesuai kebijakan operasional.

---

# 8. NON-FUNCTIONAL REQUIREMENTS

## 8.1 Performance

- Dashboard harus responsif.
- API tidak boleh memuat data berlebihan.
- Pagination digunakan pada daftar pesanan dan scan logs.
- Redirect harus ringan.
- Database query harus menggunakan index pada identifier penting.

## 8.2 Security

- Password menggunakan Argon2id atau bcrypt dengan konfigurasi yang aman.
- Semua endpoint admin membutuhkan autentikasi.
- Role dan permission diperiksa di backend.
- Input divalidasi.
- URL redirect dibatasi ke domain yang diperbolehkan.
- File upload divalidasi.
- Rate limiting diterapkan.
- Secret disimpan melalui environment variables.
- HTTPS diwajibkan saat production.
- Audit log disediakan untuk tindakan sensitif.

## 8.3 Reliability

- QR identifier harus unik.
- QR aktif harus memiliki tujuan valid.
- Perubahan URL harus bersifat konsisten.
- Produk yang sudah pensiun tidak boleh diam-diam diarahkan ke bisnis lain.
- Database memiliki backup berkala.

## 8.4 Scalability

Arsitektur MVP menggunakan modular monolith.

Modul bisnis dipisahkan secara logis:
- Authentication.
- Customers.
- Stores.
- Orders.
- Products.
- QR.
- NFC.
- Redirect.
- Analytics.
- Production.
- Shipping.

Microservices tidak diperlukan pada tahap awal.

## 8.5 Usability

- Dashboard menggunakan navigasi sederhana.
- Status pesanan mudah dibaca.
- Form menggunakan validasi.
- Operator tidak perlu memahami detail teknis backend.
- Tampilan responsif pada desktop dan tablet.

---

# 9. TECHNICAL ARCHITECTURE

## 9.1 Technology Stack

### Frontend
- Nuxt 4.
- Vue 3.
- TypeScript.
- Tailwind CSS.
- Nuxt UI.
- Pinia jika state management kompleks diperlukan.
- Zod atau validasi setara pada form jika diperlukan.

### Backend
- Go.
- Chi Router.
- REST API.
- SQLC atau SQLX.
- Structured logging.
- JWT atau session authentication.

### Database
- MySQL 8.
- Database migrations.
- Foreign key.
- Indexing.
- Transaction management.

### Supporting Tools
- Docker Compose.
- Git.
- Figma.
- NFC Tools.
- QR generation library.
- Object storage untuk asset produksi pada tahap deployment.

## 9.2 Architecture Pattern

Modular Monolith dengan pemisahan layer:

- Handler.
- Service.
- Repository.
- Domain.
- Infrastructure.

Backend menangani semua business rules.

Nuxt tidak boleh menjadi tempat penyimpanan logika bisnis kritis.

## 9.3 High-Level Architecture

Customer scan:
1. Smartphone membuka URL REPUTAP.
2. Request masuk ke Go Redirect Engine.
3. Go mencari identifier di database.
4. Sistem mengambil URL tujuan.
5. Sistem mencatat event scan.
6. Backend mengirim HTTP redirect.
7. Browser membuka Google Review.

Admin:
1. Admin membuka Nuxt dashboard.
2. Nuxt memanggil Go REST API.
3. Backend memvalidasi request.
4. Backend menjalankan business logic.
5. MySQL menyimpan data.
6. Nuxt menampilkan hasil.

## 9.4 Project Structure

```text
reputap/
│
├── apps/
│   ├── web/
│   │   ├── app/
│   │   │   ├── pages/
│   │   │   ├── components/
│   │   │   ├── layouts/
│   │   │   ├── composables/
│   │   │   ├── middleware/
│   │   │   └── types/
│   │   ├── public/
│   │   └── nuxt.config.ts
│   │
│   └── api/
│       ├── cmd/
│       │   └── server/
│       ├── internal/
│       │   ├── auth/
│       │   ├── users/
│       │   ├── customers/
│       │   ├── stores/
│       │   ├── orders/
│       │   ├── products/
│       │   ├── qrcodes/
│       │   ├── nfc/
│       │   ├── redirect/
│       │   ├── analytics/
│       │   ├── production/
│       │   ├── qualitycontrol/
│       │   ├── shipping/
│       │   └── audit/
│       ├── migrations/
│       ├── pkg/
│       └── go.mod
│
├── infrastructure/
│   ├── docker/
│   └── scripts/
│
├── docs/
│   ├── prd/
│   ├── api/
│   ├── database/
│   └── architecture/
│
├── docker-compose.yml
└── README.md
```

---

# 10. DATABASE REQUIREMENTS

Database MVP harus mendukung entitas berikut.

## 10.1 users

Data pengguna internal.

Fields:
- id.
- name.
- email.
- password_hash.
- role_id.
- status.
- created_at.
- updated_at.

## 10.2 roles

Data role internal.

Fields:
- id.
- name.
- description.

## 10.3 customers

Data pemilik bisnis.

Fields:
- id.
- customer_code.
- name.
- email.
- phone.
- address.
- city.
- province.
- notes.
- created_at.
- updated_at.

## 10.4 stores

Data bisnis.

Fields:
- id.
- customer_id.
- store_code.
- name.
- slug.
- category.
- address.
- review_url.
- status.
- created_at.
- updated_at.

## 10.5 products

Master produk REPUTAP.

Fields:
- id.
- sku.
- name.
- description.
- material.
- width.
- height.
- thickness.
- production_cost.
- selling_price.
- status.

## 10.6 orders

Data pesanan.

Fields:
- id.
- order_number.
- customer_id.
- store_id.
- order_status.
- payment_status.
- subtotal.
- discount.
- total.
- notes.
- created_by.
- created_at.
- updated_at.

## 10.7 order_items

Detail produk dalam pesanan.

Fields:
- id.
- order_id.
- product_id.
- quantity.
- unit_price.
- subtotal.

## 10.8 qr_codes

Data identifier QR.

Fields:
- id.
- store_id.
- product_instance_id.
- code.
- status.
- created_at.

## 10.9 product_instances

Data setiap unit fisik.

Fields:
- id.
- order_item_id.
- store_id.
- serial_number.
- qr_code_id.
- nfc_status.
- production_status.
- qc_status.
- created_at.

Satu unit fisik memiliki identifier permanen sendiri.

## 10.10 scan_logs

Data aktivitas redirect.

Fields:
- id.
- qr_code_id.
- source.
- scanned_at.
- device_type.
- user_agent.
- referrer.
- anonymous_identifier.

## 10.11 design_assets

Data file desain.

Fields:
- id.
- order_id.
- store_id.
- file_path.
- file_type.
- version.
- approval_status.
- uploaded_by.
- created_at.

## 10.12 production_jobs

Data pekerjaan produksi.

Fields:
- id.
- order_item_id.
- assigned_to.
- status.
- started_at.
- completed_at.
- notes.

## 10.13 quality_control

Data pengujian produk.

Fields:
- id.
- product_instance_id.
- inspected_by.
- result.
- checklist_data.
- evidence_path.
- notes.
- inspected_at.

## 10.14 shipments

Data pengiriman.

Fields:
- id.
- order_id.
- courier.
- tracking_number.
- shipping_status.
- shipped_at.
- delivered_at.

## 10.15 audit_logs

Data aktivitas internal.

Fields:
- id.
- user_id.
- action.
- entity_type.
- entity_id.
- old_values.
- new_values.
- created_at.

Catatan desain database:

- Semua identifier publik harus sulit ditebak dan unik.
- Setiap unit fisik mempunyai relasi ke toko dan pesanan.
- QR tidak boleh digunakan ulang tanpa proses reassignment.
- URL tujuan harus dapat diperbarui tanpa mengubah QR identifier.
- Relasi antarentitas menggunakan foreign key.
- Data penting tidak boleh dihapus permanen tanpa kebijakan retensi yang jelas.

---

# 11. API REQUIREMENTS

Base URL:

```text
/api/v1
```

## 11.1 Authentication

POST /auth/login

POST /auth/logout

GET /auth/me

POST /auth/refresh

## 11.2 Customers

GET /customers

POST /customers

GET /customers/{id}

PUT /customers/{id}

DELETE /customers/{id}

## 11.3 Stores

GET /stores

POST /stores

GET /stores/{id}

PUT /stores/{id}

PATCH /stores/{id}/status

PUT /stores/{id}/review-url

## 11.4 Orders

GET /orders

POST /orders

GET /orders/{id}

PUT /orders/{id}

PATCH /orders/{id}/status

GET /orders/{id}/history

## 11.5 Products

GET /products

POST /products

GET /products/{id}

PUT /products/{id}

## 11.6 QR Management

GET /qr-codes

POST /qr-codes/generate

GET /qr-codes/{id}

GET /qr-codes/{id}/download

PATCH /qr-codes/{id}/status

## 11.7 NFC

GET /products/{id}/nfc

PATCH /products/{id}/nfc-status

## 11.8 Redirect

GET /r/{code}

Endpoint redirect merupakan public endpoint dan tidak membutuhkan autentikasi.

## 11.9 Analytics

GET /analytics/overview

GET /analytics/scans

GET /analytics/stores/{id}

GET /analytics/products/{id}

## 11.10 Production

GET /production/jobs

POST /production/jobs

PATCH /production/jobs/{id}/status

## 11.11 Quality Control

POST /quality-control

GET /quality-control/{product_id}

## 11.12 Shipping

GET /shipments

POST /shipments

PATCH /shipments/{id}/status

API harus menggunakan response format konsisten, validasi input, pagination, filtering, dan error handling terstandarisasi.

---

# 12. DYNAMIC QR & NFC BUSINESS RULES

## 12.1 Permanent Identifier

Setiap unit produk memiliki identifier permanen.

Contoh:

RP-000001.

QR URL:

https://reputap.id/r/RP-000001?src=qr

NFC URL:

https://reputap.id/r/RP-000001?src=nfc

Keduanya mengarah ke bisnis yang sama.

## 12.2 Why Separate QR and NFC Parameters?

Parameter source digunakan untuk membedakan jenis interaksi.

Jika QR dan NFC menggunakan URL identik tanpa parameter source, sistem tidak dapat membedakan asal interaksi secara andal.

## 12.3 Destination Management

URL tujuan disimpan di database.

Contoh:

RP-000001 → Google Review Kopi Senja.

Jika pelanggan mengganti URL Google Review:

Admin mengubah destination URL.

QR dan NFC tetap menggunakan URL REPUTAP yang sama.

## 12.4 Product Reassignment

Jika satu display berpindah pemilik:
- Admin memverifikasi kepemilikan.
- Sistem menonaktifkan relasi sebelumnya.
- Admin menghubungkan unit ke bisnis baru.
- Riwayat lama tetap disimpan.
- QR baru hanya diaktifkan setelah konfigurasi dan QC.

## 12.5 Redirect Failure Handling

Jika URL tidak valid atau toko nonaktif:
- Sistem tidak melakukan redirect sembarangan.
- Pengunjung menerima halaman informasi.
- Admin dapat melihat error log.

---

# 13. UI/UX REQUIREMENTS

## 13.1 Design Direction

Karakter visual REPUTAP:
- Modern.
- Minimalist.
- Professional.
- Clean.
- Trustworthy.

Warna utama:
- Navy.
- Blue.
- White.
- Light Gray.

Typography:
- Poppins.

## 13.2 Main Navigation

Sidebar:

1. Overview.
2. Orders.
3. Customers.
4. Stores.
5. Products.
6. QR Management.
7. Production.
8. Quality Control.
9. Shipping.
10. Analytics.
11. Activity Logs.
12. Settings.

## 13.3 Dashboard Pages

### Overview Page
- Summary cards.
- Order chart.
- Scan chart.
- Recent orders.
- Recent activities.

### Orders Page
- Order table.
- Search.
- Filter.
- Status badge.
- Create order.
- Order detail.

### Order Detail Page
- Customer information.
- Store information.
- Product list.
- QR configuration.
- NFC configuration.
- Design preview.
- Production timeline.
- QC results.
- Shipping information.

### QR Management Page
- QR list.
- QR preview.
- Download QR.
- Status.
- Destination URL.
- Source tracking.

### Production Page
- Production queue.
- Assigned operator.
- Production checklist.
- Notes.

### Analytics Page
- Total scans.
- QR vs NFC.
- Daily scans.
- Store comparison.
- Product comparison.

---

# 14. SECURITY & PRIVACY

## 14.1 Authentication Security

- Password hashing.
- Rate limiting.
- Session/token expiry.
- Permission validation.
- Secure cookies jika menggunakan cookie-based authentication.
- CSRF protection bila menggunakan cookie authentication.

## 14.2 Redirect Security

Redirect endpoint harus membatasi URL tujuan ke domain Google yang telah ditentukan.

Sistem tidak boleh mengizinkan admin atau pihak luar menjadikan REPUTAP sebagai open redirect menuju sembarang website.

URL harus diperiksa saat disimpan dan tetap divalidasi ketika redirect dilakukan.

## 14.3 Scan Data Privacy

Sistem hanya mengumpulkan data yang dibutuhkan untuk analitik.

- Jangan mengumpulkan informasi personal yang tidak diperlukan.
- Hindari penyimpanan IP permanen pada MVP.
- Gunakan identifier anonim secara proporsional.
- Tetapkan kebijakan retensi scan logs.
- Sediakan akses terbatas untuk data operasional.

## 14.4 Google Reviews Policy

REPUTAP hanya membantu pelanggan membuka halaman review.

Sistem dan tim tidak boleh:
- Membeli ulasan.
- Menghasilkan ulasan palsu.
- Menjanjikan rating tertentu.
- Memberikan insentif yang bertujuan memengaruhi ulasan positif.
- Mengarahkan hanya pelanggan puas ke Google Review.

Semua pelanggan harus memiliki kesempatan yang wajar untuk memberikan ulasan secara jujur dan sukarela.

---

# 15. ACCEPTANCE CRITERIA

MVP dinyatakan memenuhi kebutuhan inti apabila:

## AC-001 — Authentication
Admin dapat login dan logout.

## AC-002 — Customer Management
Admin dapat menambahkan dan mengelola pelanggan.

## AC-003 — Store Management
Admin dapat mendaftarkan toko dan menyimpan Google Review URL.

## AC-004 — Order Management
Admin dapat membuat pesanan dan memperbarui statusnya.

## AC-005 — Dynamic QR
Sistem dapat membuat QR unik berdasarkan identifier produk.

## AC-006 — Redirect
QR yang dipindai berhasil mengarahkan pengunjung ke Google Review yang benar.

## AC-007 — Dynamic Destination
Perubahan URL tujuan tidak mengubah QR fisik.

## AC-008 — NFC
NFC dapat ditulis menggunakan URL REPUTAP dan berhasil mengarahkan pengunjung ke tujuan yang benar.

## AC-009 — Source Tracking
Sistem dapat membedakan QR scan dan NFC tap.

## AC-010 — Design Management
Admin dapat mengunggah dan menyimpan desain produk.

## AC-011 — Production Workflow
Operator dapat memperbarui tahapan produksi.

## AC-012 — Quality Control
Produk yang gagal QC tidak dapat ditandai siap dikirim.

## AC-013 — Analytics
Dashboard dapat menampilkan jumlah scan dan statistik penggunaan.

## AC-014 — Audit Trail
Perubahan penting tercatat.

## AC-015 — Security
Endpoint admin terlindungi dan redirect tidak dapat digunakan sebagai open redirect.

---

# 16. MVP SCOPE & PRIORITIZATION

## P0 — Critical / Must Have

- Admin login.
- Dashboard overview.
- Customer management.
- Store management.
- Order management.
- Product management.
- QR generation.
- Dynamic redirect.
- Google Review URL management.
- NFC URL configuration.
- Production status.
- QC checklist.
- Basic scan analytics.

## P1 — Important / Should Have

- Design asset management.
- Design revision.
- Shipping management.
- Operator assignment.
- Advanced filtering.
- Activity logs.
- QR/NFC source analytics.
- Export data.

## P2 — Future Development

- Customer portal.
- Multi-tenant SaaS.
- Subscription.
- Payment gateway.
- WhatsApp notification automation.
- Advanced analytics.
- Inventory management.
- Automated vendor integration.
- Multi-branch business management.
- White-label solution.

---

# 17. DEVELOPMENT ROADMAP

## PHASE 1 — Foundation

Target: project dapat dijalankan secara lokal.

Tasks:
- Setup repository.
- Setup Nuxt.
- Setup Go.
- Setup MySQL.
- Setup Docker Compose.
- Setup migrations.
- Setup API conventions.
- Setup authentication.

Deliverable:
Project skeleton siap dikembangkan.

## PHASE 2 — Core Business Management

Tasks:
- Customer CRUD.
- Store CRUD.
- Product CRUD.
- Order CRUD.
- Order status workflow.

Deliverable:
Tim dapat mencatat dan mengelola pesanan.

## PHASE 3 — QR & NFC Engine

Tasks:
- Product identifier generation.
- QR generator.
- QR download.
- NFC URL generator.
- Redirect endpoint.
- Destination update.
- Source tracking.

Deliverable:
QR dan NFC dapat diuji tanpa produk fisik.

## PHASE 4 — Internal Operations

Tasks:
- Design management.
- Production workflow.
- QC management.
- Shipping management.
- Activity logs.

Deliverable:
Tim dapat menangani pesanan dari konfigurasi sampai pengiriman.

## PHASE 5 — Analytics & Testing

Tasks:
- Scan analytics.
- QR vs NFC reporting.
- Performance testing.
- Security testing.
- Integration testing.
- End-to-end testing.

Deliverable:
MVP siap digunakan untuk uji pasar.

## PHASE 6 — Physical Product Integration

Tasks:
- Produksi prototype.
- Encode NFC.
- Cetak QR.
- Pemasangan stiker.
- QC fisik.
- Pengujian dengan smartphone.
- Uji coba pada UMKM.

Deliverable:
Produk fisik dan sistem digital berjalan sebagai satu kesatuan.

---

# 18. TESTING STRATEGY

## 18.1 Unit Testing

Test:
- Identifier generation.
- URL validation.
- QR generation.
- Redirect resolution.
- Order status transition.
- Permission validation.

## 18.2 Integration Testing

Test:
- Order ke product instance.
- Product instance ke QR.
- QR ke store.
- Store ke destination URL.
- NFC URL ke redirect engine.
- Scan log ke analytics.

## 18.3 End-to-End Testing

Skenario utama:

1. Admin login.
2. Membuat customer.
3. Membuat store.
4. Memasukkan Google Review URL.
5. Membuat order.
6. Menghasilkan product instance.
7. Generate QR.
8. Membuka QR melalui smartphone.
9. Memastikan redirect berhasil.
10. Mengubah Google Review URL.
11. Membuka QR yang sama.
12. Memastikan tujuan sudah berubah.
13. Menguji NFC.
14. Memeriksa scan analytics.

## 18.4 Edge Case Testing

- QR tidak ditemukan.
- QR nonaktif.
- Store nonaktif.
- Google Review URL tidak valid.
- NFC gagal dibaca.
- Product instance belum memiliki destination.
- Pesanan dibatalkan.
- QC gagal.
- Redirect database timeout.
- Duplicate order.
- QR identifier duplikat.

---

# 19. ASSUMPTIONS & CONSTRAINTS

## 19.1 Assumptions

- Tim REPUTAP menangani seluruh konfigurasi produk.
- Customer tidak memerlukan dashboard.
- Google Business Profile telah dimiliki atau dikelola oleh pelanggan.
- Pelanggan menyediakan URL Google Review yang benar.
- Setiap unit produk memiliki QR identifier permanen.
- Setiap unit fisik memiliki satu konfigurasi tujuan utama.
- NFC dan QR menggunakan redirect engine yang sama.

## 19.2 Constraints

- MVP dikembangkan dengan modal terbatas.
- Produksi fisik dilakukan dalam batch kecil.
- Sistem pertama kali berjalan di local environment.
- Integrasi pembayaran dan kurir belum menjadi kebutuhan awal.
- Google tidak menjamin ulasan langsung dipublikasikan.
- Sistem tidak dapat mengukur jumlah review yang benar-benar diterbitkan hanya dari scan event.

## 19.3 Important Technical Constraint

Dynamic QR membutuhkan domain REPUTAP yang dapat diakses publik saat produk mulai digunakan.

QR berbasis localhost hanya untuk testing.

Jika domain REPUTAP tidak dapat diakses, pengunjung tidak dapat menggunakan fitur redirect.

---

# 20. RISKS & MITIGATION

| Risk | Impact | Mitigation |
|---|---|---|
| Salah URL toko | Pengunjung masuk ke bisnis lain | URL verification dan QC |
| NFC gagal dibaca | Produk tidak berfungsi optimal | NFC testing |
| QR rusak saat dicetak | Pengunjung gagal scan | Print testing |
| Domain tidak tersedia | Semua redirect terganggu | Monitoring dan backup |
| Salah konfigurasi produk | Kerugian produksi | Product identifier dan checklist |
| QR disalahgunakan | Penyalahgunaan redirect | Allowlist destination |
| Data pelanggan bocor | Risiko privasi | Access control dan encryption |
| Produk rusak saat pengiriman | Komplain pelanggan | Packaging dan QC |
| Ulasan tidak dipublikasikan | Ekspektasi pelanggan tidak sesuai | Edukasi bahwa publikasi mengikuti kebijakan Google |

---

# 21. OUT OF SCOPE — MVP

Fitur berikut tidak dikerjakan pada tahap pertama:

- Customer login.
- Customer self-service dashboard.
- Subscription billing.
- Payment gateway.
- Google Reviews API integration untuk mengambil seluruh ulasan.
- Automatic review generation.
- Rating manipulation.
- AI sentiment analysis.
- CRM marketing automation.
- WhatsApp chatbot.
- Mobile application native.
- Multi-tenant SaaS billing.

---

# 22. FINAL PRODUCT DEFINITION

REPUTAP adalah layanan produk fisik berbasis akrilik yang mengintegrasikan NFC dan Dynamic QR untuk mempermudah pelanggan mengakses halaman Google Review suatu bisnis.

Sistem internal REPUTAP digunakan oleh tim untuk mengelola:
- Pelanggan.
- Bisnis.
- Pesanan.
- QR Code.
- NFC.
- URL tujuan.
- Desain.
- Produksi.
- Quality Control.
- Pengiriman.
- Statistik penggunaan.

Pelanggan tidak diwajibkan memahami sistem teknis maupun mengelola dashboard.

Tim REPUTAP bertanggung jawab menyiapkan produk sampai siap digunakan.

## FINAL PRODUCT PRINCIPLE

**REPUTAP handles the technology. Businesses focus on their customers.**

Produk harus sesederhana mungkin bagi pemilik bisnis, tetapi memiliki sistem operasional yang terstruktur di belakangnya.

---

**END OF DOCUMENT**

**REPUTAP PRD v2.0 — Done-for-You Smart Review Platform**