CREATE TABLE product_categories(
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(255)    NOT NULL,
    description     TEXT            NULL,
    business_id     UUID            NOT NULL,
    created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    CONSTRAINT fk_business
        FOREIGN KEY (business_id) REFERENCES business(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_categories_business_id ON product_categories(business_id);
CREATE INDEX idx_product_categories_name ON product_categories(name);
