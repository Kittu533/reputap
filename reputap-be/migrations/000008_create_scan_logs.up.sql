CREATE TABLE IF NOT EXISTS scan_logs (
    id BIGSERIAL PRIMARY KEY,
    qr_code_id BIGINT NOT NULL,
    source VARCHAR(50) DEFAULT 'unknown',
    device_type VARCHAR(50) DEFAULT 'other',
    user_agent TEXT,
    referrer VARCHAR(500),
    scanned_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (qr_code_id) REFERENCES qr_codes(id)
);

CREATE INDEX IF NOT EXISTS idx_scan_logs_qr_code ON scan_logs (qr_code_id);
CREATE INDEX IF NOT EXISTS idx_scan_logs_scanned_at ON scan_logs (scanned_at);
CREATE INDEX IF NOT EXISTS idx_scan_logs_source ON scan_logs (source);
