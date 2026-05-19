-- +migrate Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE TYPE IF NOT EXISTS group_type AS ENUM {
    'GROUP',
    'CLUB'
};

CREATE TABLE IF NOT EXISTS "groups" (
    uuid           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    "name"         VARCHAR(45) NOT NULL,
    "description"  TEXT        NULL,
    "group_type"   group_type  NOT NULL,
    "parent_group_uuid" UUID        NULL REFERENCES groups(uuid) ON UPDATE CASCADE ON DELETE SET NULL,

    deleted_at     TIMESTAMPTZ NULL,
    updated_at     TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    created_at     TIMESTAMPTZ DEFAULT NOW() NOT NULL
); --да, я люблю структурировать бд
CREATE INDEX idx_groups_name_trgm ON groups USING GIN ("name" gin_trgm_ops);

CREATE TABLE IF NOT EXISTS "groups_members" (
    group_uuid  UUID NOT NULL REFERENCES groups(uuid) ON UPDATE CASCADE ON DELETE CASCADE,
    member_uuid UUID NOT NULL, --SAGA
    joined_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    PRIMARY KEY(group_uuid, member_uuid)
);
CREATE INDEX idx_groups_members_member_uuid ON "groups_members" ("member_uuid");

CREATE TABLE IF NOT EXISTS "groups_curators" (
    group_uuid   UUID NOT NULL REFERENCES groups(uuid) ON UPDATE CASCADE ON DELETE CASCADE,
    curator_uuid UUID NOT NULL, --SAGA
    assigned_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    PRIMARY KEY(group_uuid, curator_uuid)
);
CREATE INDEX idx_groups_curators_group_uuid ON "groups_curators" ("group_uuid");

CREATE TABLE IF NOT EXISTS "curators_hierarchy" (
    senior_uuid      UUID NOT NULL,
    subordinate_uuid UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    PRIMARY KEY(senior_uuid, subordinate_uuid),
    CHECK senior_uuid!=subordinate_uuid --так тоже можно
);
CREATE INDEX idx_curators_hierarchy_subordinate_uuid ON "curators_hierarchy" ("subordinate_uuid");
