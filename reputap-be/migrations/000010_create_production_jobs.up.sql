CREATE TABLE IF NOT EXISTS production_jobs (
    id BIGSERIAL PRIMARY KEY,
    product_instance_id BIGINT NOT NULL,
    assigned_to BIGINT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    checklist_data JSONB,
    started_at TIMESTAMPTZ NULL,
    completed_at TIMESTAMPTZ NULL,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_instance_id) REFERENCES product_instances(id),
    FOREIGN KEY (assigned_to) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_production_jobs_status ON production_jobs (status);
CREATE INDEX IF NOT EXISTS idx_production_jobs_assigned ON production_jobs (assigned_to);
