CREATE TABLE IF NOT EXISTS quality_controls (
    id BIGSERIAL PRIMARY KEY,
    product_instance_id BIGINT NOT NULL,
    inspected_by BIGINT NOT NULL,
    result VARCHAR(50) NOT NULL,
    checklist_data JSONB NOT NULL,
    evidence_paths JSONB,
    notes TEXT,
    inspected_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_instance_id) REFERENCES product_instances(id),
    FOREIGN KEY (inspected_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_qc_product_instance ON quality_controls (product_instance_id);
