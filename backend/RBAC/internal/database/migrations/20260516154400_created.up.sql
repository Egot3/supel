-- +migrate Up
CREATE TYPE IF NOT EXISTS verb AS ENUM (
    'GET', 
    'POST', 
    'PUT',
    'PATCH',
    'DELETE'
);

CREATE TABLE IF NOT EXISTS "actions" (
    action_uuid UUID              PRIMARY KEY DEFAULT gen_random_uuid(),
    scope       character varying NOT NULL,
    sub_scope   character varying NULL,
    "verb"      verb              NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_actions_full ON "actions" ("scope", "sub_scope", "verb") WHERE "sub_scope" IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_actions_subscope_null ON "actions" ("scope", "verb") WHERE "sub_scope" IS NULL;

CREATE TABLE IF NOT EXISTS "roles" (
    role_uuid        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name        character varying NOT NULL UNIQUE,
    role_description TEXT NULL,

    "priority" smallint DEFAULT 0 NOT NULL, 

    extended_role_uuid     UUID NULL REFERENCES roles(role_uuid) ON UPDATE CASCADE ON DELETE SET NULL,

    created_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_role_uuid_active ON roles (role_uuid) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS "roles_actions" (
    action_uuid UUID REFERENCES actions(action_uuid) ON DELETE CASCADE ON UPDATE CASCADE,
    role_uuid   UUID REFERENCES roles(role_uuid) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY(action_uuid, role_uuid)
);
CREATE INDEX IF NOT EXISTS idx_role_uuid_junction ON "roles_actions" ("role_uuid");

CREATE TABLE IF NOT EXISTS "users_roles" (
    user_uuid     UUID NOT NULL, --SAGA
    role_uuid     UUID REFERENCES roles(role_uuid) ON UPDATE CASCADE ON DELETE CASCADE,
    assignor_uuid UUID NOT NULL, --SAGA 

    assigned_at TIMESTAMPTZ DEFAULT now(),
    expires_at  TIMESTAMPTZ NULL,

    PRIMARY KEY(user_uuid, role_uuid)
);
CREATE INDEX IF NOT EXISTS idx_user_uuid ON "users_roles" ("user_uuid");