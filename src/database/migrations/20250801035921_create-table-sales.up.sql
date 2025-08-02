CREATE TABLE sales (
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    outlet_id       UUID NOT NULL,
    outlet_staff_id UUID NOT NULL,
    customer_id     UUID NULL,
    payment_method_id UUID NULL,
    table_id        UUID NOT NULL,
    invoice_number  VARCHAR(50) NOT NULL UNIQUE,
    total           NUMERIC(10, 2) NOT NULL,
    discount        NUMERIC(10, 2) DEFAULT 0 NOT NULL,
    tax             NUMERIC(10, 2) DEFAULT 0 NOT NULL,
    grand_total     NUMERIC(10, 2) NOT NULL,
    status          VARCHAR(50) NOT NULL, -- paid, unpaid, void, hold
    sale_date       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    note            TEXT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_outlet
        FOREIGN KEY (outlet_id) REFERENCES outlets(id) ON DELETE CASCADE,
    CONSTRAINT fk_table
        FOREIGN KEY (table_id) REFERENCES tables(id) ON DELETE CASCADE,
    CONSTRAINT fk_outlet_staff
        FOREIGN KEY (outlet_staff_id) REFERENCES outlet_staff(id) ON DELETE CASCADE,
    CONSTRAINT fk_customer
        FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL,  
    CONSTRAINT fk_payment_method
        FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id) ON DELETE SET NULL
 
);

CREATE INDEX idx_sales_outlet_id ON sales(outlet_id);
CREATE INDEX idx_sales_outlet_staff_id ON sales(outlet_staff_id);
CREATE INDEX idx_sales_customer_id ON sales(customer_id);
CREATE INDEX idx_sales_payment_method_id ON sales(payment_method_id);
CREATE INDEX idx_sales_table_id ON sales(table_id);
CREATE INDEX idx_sales_invoice_number ON sales(invoice_number);
CREATE INDEX idx_sales_sale_date ON sales(sale_date);