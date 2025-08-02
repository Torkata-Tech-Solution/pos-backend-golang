CREATE TABLE sales_items (
     id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    sale_id         UUID NOT NULL,
    product_id      UUID NOT NULL,
    quantity        INT NOT NULL,
    price           NUMERIC(10, 2) NOT NULL,
    discount        NUMERIC(10, 2) DEFAULT 0 NOT NULL,
    total           NUMERIC(10, 2) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sale
        FOREIGN KEY (sale_id) REFERENCES sales(id) ON DELETE CASCADE,
    CONSTRAINT fk_product
        FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_sales_items_sale_id ON sales_items(sale_id);
CREATE INDEX idx_sales_items_product_id ON sales_items(product_id);