CREATE TABLE scan_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    qr_code_id BIGINT UNSIGNED NOT NULL,
    source ENUM('qr', 'nfc', 'unknown') DEFAULT 'unknown',
    device_type ENUM('android', 'ios', 'desktop', 'other') DEFAULT 'other',
    user_agent TEXT,
    referrer VARCHAR(500),
    scanned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (qr_code_id) REFERENCES qr_codes(id),
    INDEX idx_scan_logs_qr_code (qr_code_id),
    INDEX idx_scan_logs_scanned_at (scanned_at),
    INDEX idx_scan_logs_source (source)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
