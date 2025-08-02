CREATE TABLE business_users(
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    business_id     UUID            NOT NULL,
    user_id         UUID            NOT NULL,
    role            VARCHAR(255)    NOT NULL,
    created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    CONSTRAINT fk_business
        FOREIGN KEY (business_id) REFERENCES business(id) ON DELETE CASCADE,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_business_users_business_id ON business_users(business_id);
