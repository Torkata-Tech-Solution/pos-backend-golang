CREATE TABLE tables (
     id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    outlet_id   UUID            NOT NULL,
    name        VARCHAR(255)    NOT NULL, -- Nama atau nomor meja, misal: "A1"
    location    VARCHAR(255), -- Opsional, misal: "Lantai 2"
    status      VARCHAR(50), -- available, occupied, reserved
    capacity    INTEGER         NOT NULL, -- Jumlah maksimum pelanggan
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_outlet
        FOREIGN KEY (outlet_id) REFERENCES outlets(id) ON DELETE CASCADE
);

CREATE INDEX idx_tables_outlet_id ON tables(outlet_id);
CREATE INDEX idx_tables_name ON tables(name);