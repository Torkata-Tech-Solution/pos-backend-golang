CREATE TABLE payment_methods (
     id              UUID     PRIMARY KEY DEFAULT uuid_generate_v4(),
    outlet_id   UUID NOT NULL,
    name        VARCHAR(255) NOT NULL,
    type        VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_outlet
        FOREIGN KEY (outlet_id) REFERENCES outlets(id) ON DELETE CASCADE
);

CREATE INDEX idx_payment_methods_outlet_id ON payment_methods(outlet_id);
CREATE INDEX idx_payment_methods_name ON payment_methods(name);