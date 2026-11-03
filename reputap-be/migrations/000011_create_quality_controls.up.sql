CREATE TABLE quality_controls (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_instance_id BIGINT UNSIGNED NOT NULL,
    inspected_by BIGINT UNSIGNED NOT NULL,
    result ENUM('passed', 'failed', 'rework_required') NOT NULL,
    checklist_data JSON NOT NULL,
    evidence_paths JSON,
    notes TEXT,
    inspected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_instance_id) REFERENCES product_instances(id),
    FOREIGN KEY (inspected_by) REFERENCES users(id),
    INDEX idx_qc_product_instance (product_instance_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
