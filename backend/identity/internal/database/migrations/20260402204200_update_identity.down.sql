-- +migrate Down
ALTER TABLE users 
    DROP IF EXISTS COLUMN password_hash,
    DROP IF EXISTS COLUMN email;
