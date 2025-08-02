CREATE TABLE outlet_staff(
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    outlet_id       UUID            NOT NULL,
    name            VARCHAR(255)    NOT NULL,
    username        VARCHAR(255)    NULL UNIQUE,
    password        VARCHAR(255)    NOT NULL,
    role            VARCHAR(255)    NOT NULL,
    created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    CONSTRAINT fk_outlet
        FOREIGN KEY (outlet_id) REFERENCES outlets(id) ON DELETE CASCADE
);
CREATE INDEX idx_outlet_staff_outlet_id ON outlet_staff(outlet_id);
