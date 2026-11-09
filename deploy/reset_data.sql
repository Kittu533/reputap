-- ====================================================================
-- REPUTAP: Clean & Reset Database for Production
-- Menghapus semua data dummy transaksi, customer, store, order, QR, scan
-- Menyisakan: roles, permissions, products, dan user Super Admin
-- ====================================================================

BEGIN;

-- 1. Kosongkan semua tabel operasional & transaksi dan reset ID kembali ke 1
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

-- 2. Hapus demo users selain super admin (admin@reputap.id)
DELETE FROM users WHERE email != 'admin@reputap.id';

-- 3. Pastikan user super admin aktif
INSERT INTO users (name, email, password_hash, role_id, status)
VALUES ('Super Admin', 'admin@reputap.id', '$2a$10$BWefn4Lz1nFfOA92nVSqFuOhUPCI/RXFxTEBz0l6G6hOxHPEBli7e', 1, 'active')
ON CONFLICT (email) DO UPDATE SET status = 'active', role_id = 1;

COMMIT;
