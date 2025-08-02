CREATE TABLE products(
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    image           VARCHAR(255)    NULL,
    name            VARCHAR(255)    NOT NULL,
    description     TEXT            NULL,
    price           DECIMAL(10, 2)  NOT NULL,
    category_id     UUID            NOT NULL,
    business_id     UUID            NOT NULL,
    created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    CONSTRAINT fk_business
        FOREIGN KEY (business_id) REFERENCES business(id) ON DELETE CASCADE,
    CONSTRAINT fk_category
        FOREIGN KEY (category_id) REFERENCES product_categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_products_business_id ON products(business_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_name ON products(name);
