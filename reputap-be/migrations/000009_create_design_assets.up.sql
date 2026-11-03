CREATE TABLE design_assets (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    store_id BIGINT UNSIGNED NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(10) NOT NULL,
    file_size BIGINT UNSIGNED DEFAULT 0,
    version INT UNSIGNED DEFAULT 1,
    approval_status ENUM('draft', 'pending_approval', 'approved', 'revision_required') DEFAULT 'draft',
    uploaded_by BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (store_id) REFERENCES stores(id),
    FOREIGN KEY (uploaded_by) REFERENCES users(id),
    INDEX idx_design_assets_order (order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
