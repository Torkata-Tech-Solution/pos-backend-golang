CREATE TABLE coupons (
     id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    outlet_id       UUID NOT NULL,
    code            VARCHAR(100) NOT NULL UNIQUE,
    description     TEXT NULL,
    discount_type   VARCHAR(50) NOT NULL,
    discount_value  NUMERIC(10, 2) NOT NULL,
    max_uses        INT NOT NULL DEFAULT 1,
    used_count      INT NOT NULL DEFAULT 0,
    start_date      TIMESTAMP NOT NULL,
    end_date        TIMESTAMP NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_outlet
        FOREIGN KEY (outlet_id) REFERENCES outlets(id) ON DELETE CASCADE
);

CREATE INDEX idx_coupons_outlet_id ON coupons(outlet_id);
CREATE INDEX idx_coupons_code ON coupons(code);