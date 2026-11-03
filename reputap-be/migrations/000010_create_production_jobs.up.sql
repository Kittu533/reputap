CREATE TABLE production_jobs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_instance_id BIGINT UNSIGNED NOT NULL,
    assigned_to BIGINT UNSIGNED NULL,
    status ENUM('pending', 'in_progress', 'completed', 'failed') DEFAULT 'pending',
    checklist_data JSON,
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_instance_id) REFERENCES product_instances(id),
    FOREIGN KEY (assigned_to) REFERENCES users(id),
    INDEX idx_production_jobs_status (status),
    INDEX idx_production_jobs_assigned (assigned_to)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
