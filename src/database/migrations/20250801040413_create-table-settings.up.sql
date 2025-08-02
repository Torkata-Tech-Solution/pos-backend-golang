CREATE TABLE settings (
     id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    outlet_id       UUID NOT NULL,
    key             VARCHAR(100) NOT NULL,
    value           TEXT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_outlet
        FOREIGN KEY (outlet_id) REFERENCES outlets(id) ON DELETE CASCADE
);

CREATE INDEX idx_settings_outlet_id ON settings(outlet_id);
CREATE INDEX idx_settings_key ON settings(key);