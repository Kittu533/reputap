CREATE TABLE product_instances (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_item_id BIGINT UNSIGNED NOT NULL,
    store_id BIGINT UNSIGNED NOT NULL,
    serial_number VARCHAR(20) NOT NULL UNIQUE,
    nfc_status ENUM('not_configured', 'configured', 'tested', 'failed', 'replaced') DEFAULT 'not_configured',
    nfc_configured_at TIMESTAMP NULL,
    nfc_configured_by BIGINT UNSIGNED NULL,
    production_status ENUM('pending', 'in_progress', 'completed', 'failed') DEFAULT 'pending',
    qc_status ENUM('pending', 'passed', 'failed', 'rework') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (order_item_id) REFERENCES order_items(id),
    FOREIGN KEY (store_id) REFERENCES stores(id),
    FOREIGN KEY (nfc_configured_by) REFERENCES users(id),
    UNIQUE INDEX idx_product_instances_serial (serial_number),
    INDEX idx_product_instances_order_item (order_item_id),
    INDEX idx_product_instances_store (store_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE qr_codes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_instance_id BIGINT UNSIGNED NOT NULL UNIQUE,
    store_id BIGINT UNSIGNED NOT NULL,
    code VARCHAR(20) NOT NULL UNIQUE,
    status ENUM('active', 'inactive', 'suspended', 'retired') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_instance_id) REFERENCES product_instances(id),
    FOREIGN KEY (store_id) REFERENCES stores(id),
    UNIQUE INDEX idx_qr_codes_code (code),
    INDEX idx_qr_codes_status (status),
    INDEX idx_qr_codes_store (store_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
