CREATE TABLE IF NOT EXISTS shipments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    recipient_name VARCHAR(150) NOT NULL,
    recipient_address TEXT NOT NULL,
    recipient_phone VARCHAR(20),
    courier VARCHAR(100),
    tracking_number VARCHAR(100),
    shipping_status VARCHAR(50) DEFAULT 'waiting',
    shipped_at TIMESTAMPTZ NULL,
    delivered_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE INDEX IF NOT EXISTS idx_shipments_order ON shipments (order_id);
CREATE INDEX IF NOT EXISTS idx_shipments_status ON shipments (shipping_status);
