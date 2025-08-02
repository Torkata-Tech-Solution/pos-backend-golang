CREATE TABLE printers (
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    outlet_id       UUID NOT NULL,
    name            VARCHAR(255) NOT NULL,
    connection_type VARCHAR(50) NOT NULL,
    mac_address     VARCHAR(50) NULL,
    ip_address      VARCHAR(50) NULL,
    paper_width     INTEGER NULL,
    default_printer BOOLEAN DEFAULT FALSE NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_outlet
        FOREIGN KEY (outlet_id) REFERENCES outlets(id) ON DELETE CASCADE
);

CREATE INDEX idx_printers_outlet_id ON printers(outlet_id);
CREATE INDEX idx_printers_name ON printers(name);