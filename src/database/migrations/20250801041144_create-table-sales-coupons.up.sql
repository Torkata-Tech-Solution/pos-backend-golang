CREATE TABLE sales_coupons (
     id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    sale_id         UUID NOT NULL,
    coupon_id       UUID NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sale
        FOREIGN KEY (sale_id) REFERENCES sales(id) ON DELETE CASCADE,
    CONSTRAINT fk_coupon
        FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE CASCADE
);

CREATE INDEX idx_sales_coupons_sale_id ON sales_coupons(sale_id);
CREATE INDEX idx_sales_coupons_coupon_id ON sales_coupons(coupon_id);