-- +migrate Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- CREATE INDEX idx_groups_name_trgm ON groups USING GIN ("name" gin_trgm_ops);

CREATE TYPE IF NOT EXISTS puddle_type AS ENUM {
    'GROUP',
    'ONEONONE',
    'CHANNEL'
};

CREATE TABLE IF NOT EXISTS "puddles" {
    uuid                UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    "name"              VARCHAR(45)  NOT NULL,
    "description"       TEXT         NULL,
    "puddle_type"       puddle_type  NOT NULL,

    deleted_at     TIMESTAMPTZ NULL,
    updated_at     TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    created_at     TIMESTAMPTZ DEFAULT NOW() NOT NULL
};
CREATE INDEX idx_puddles_name_trgm ON puddles USING GIN ("name" gin_trgm_ops);

CREATE TABLE IF NOT EXISTS "puddles_members" {
    puddle_uuid     UUID REFERENCES puddles(uuid) ON UPDATE CASCADE ON DELETE CASCADE,
    member_uuid     UUID NOT NULL,
    adder_uuid      UUID NOT NULL,
    joined_at       TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY(puddle_uuid, member_uuid)
};
CREATE INDEX idx_puddles_members_member_uuid ON "puddles_members" ("member_uuid");
CREATE INDEX idx_puddles_members_puddle_uuid ON "groups_members" ("member_uuid");

CREATE TABLE IF NOT EXISTS "puddles_moderators" {
    puddle_uuid    UUID        REFERENCES puddles(uuid) ON UPDATE CASCADE ON DELETE CASCADE,
    moderator_uuid UUID        NOT NULL,
    assignor_uuid  UUID        NOT NULL,
    assigned_at    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    PRIMARY KEY(puddle_uuid, moderator_uuid)
};
CREATE INDEX idx_puddles_moderators_puddle_uuid ON "puddles_moderators" ("puddle_uuid");
