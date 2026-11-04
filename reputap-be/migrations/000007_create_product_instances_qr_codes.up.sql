CREATE TABLE IF NOT EXISTS product_instances (
    id BIGSERIAL PRIMARY KEY,
    order_item_id BIGINT NULL,
    store_id BIGINT NULL,
    serial_number VARCHAR(20) NOT NULL UNIQUE,
    nfc_status VARCHAR(50) DEFAULT 'not_configured',
    nfc_configured_at TIMESTAMPTZ NULL,
    nfc_configured_by BIGINT NULL,
    production_status VARCHAR(50) DEFAULT 'pending',
    qc_status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_item_id) REFERENCES order_items(id),
    FOREIGN KEY (store_id) REFERENCES stores(id),
    FOREIGN KEY (nfc_configured_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_product_instances_order_item ON product_instances (order_item_id);
CREATE INDEX IF NOT EXISTS idx_product_instances_store ON product_instances (store_id);

CREATE TABLE IF NOT EXISTS qr_codes (
    id BIGSERIAL PRIMARY KEY,
    product_instance_id BIGINT NOT NULL UNIQUE,
    store_id BIGINT NULL,
    code VARCHAR(20) NOT NULL UNIQUE,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_instance_id) REFERENCES product_instances(id),
    FOREIGN KEY (store_id) REFERENCES stores(id)
);

CREATE INDEX IF NOT EXISTS idx_qr_codes_status ON qr_codes (status);
CREATE INDEX IF NOT EXISTS idx_qr_codes_store ON qr_codes (store_id);
