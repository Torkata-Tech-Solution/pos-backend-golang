CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users(
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    photo           VARCHAR(255)    DEFAULT NULL,
    name            VARCHAR(255)    NOT NULL,
    email           VARCHAR(255)    NOT NULL UNIQUE,
    phone           VARCHAR(20)     NOT NULL,
    password        VARCHAR(255)    NOT NULL,
    role            VARCHAR(255)    NOT NULL,
    verified_email  BOOLEAN         DEFAULT FALSE  NOT NULL,
    created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    updated_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL
);
