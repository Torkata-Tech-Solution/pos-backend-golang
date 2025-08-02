CREATE TABLE outlets(
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    business_id     UUID            NOT NULL,
    name            VARCHAR(255)    NOT NULL,
    address         VARCHAR(255)    NOT NULL,
    phone           VARCHAR(20)     NULL UNIQUE,
    email           VARCHAR(255)    NULL UNIQUE,
    created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    CONSTRAINT fk_business
        FOREIGN KEY (business_id) REFERENCES business(id) ON DELETE CASCADE
);

CREATE INDEX idx_outlet_business_id ON outlets(business_id);