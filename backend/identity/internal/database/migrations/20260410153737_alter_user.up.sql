-- Modify "users" table
ALTER IF EXISTS TABLE users ALTER COLUMN "role" TYPE character varying USING "role"::text, ALTER COLUMN "role" SET DEFAULT 'USER', ALTER COLUMN "email" TYPE character varying, ALTER COLUMN "password_hash" TYPE character varying;
-- Drop enum type "user_role"
DROP TYPE "user_role";
