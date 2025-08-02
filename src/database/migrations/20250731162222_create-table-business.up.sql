CREATE TABLE business(
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    domain          VARCHAR(255)    NOT NULL UNIQUE,
    name            VARCHAR(255)    NOT NULL,
    address         VARCHAR(255)    NOT NULL,
    phone           VARCHAR(20)     NULL UNIQUE,
    email           VARCHAR(255)    NULL UNIQUE,
    website         VARCHAR(255)    NULL,
    logo            VARCHAR(255)    NULL,
    created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL
);
CREATE INDEX idx_business_domain ON business(domain);