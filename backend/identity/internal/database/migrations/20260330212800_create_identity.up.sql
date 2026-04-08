-- +migrate Up
CREATE TYPE user_role AS ENUM('USER','ADMIN');

CREATE TABLE IF NOT EXISTS users (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role user_role NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);