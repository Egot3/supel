-- +migrate Up
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS "groups" (
    uuid                TEXT        PRIMARY KEY DEFAULT (lower(
        hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-' || '4' || substr(hex(randomblob(2)),2) || '-' || 
        substr('89ab', (abs(random()) % 4) + 1, 1) || substr(hex(randomblob(2)),2) || '-' || 
        hex(randomblob(6))
    )),
    "name"              TEXT        NOT NULL,
    "description"       TEXT,
    "group_type"        TEXT        NOT NULL CHECK (group_type IN ('GROUP', 'CLUB')),

    deleted_at          TEXT,
    updated_at          TEXT        NOT NULL DEFAULT (datetime('now')),
    created_at          TEXT        NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX idx_groups_name ON groups("name");

CREATE TABLE IF NOT EXISTS "groups_members" (
    group_uuid  TEXT NOT NULL REFERENCES groups(uuid) ON UPDATE CASCADE ON DELETE CASCADE,
    member_uuid TEXT NOT NULL, --SAGA
    joined_at   TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (group_uuid, member_uuid)
);
CREATE INDEX idx_groups_members_member_uuid ON "groups_members" ("member_uuid");

CREATE TABLE IF NOT EXISTS "groups_curators" (
    group_uuid   TEXT NOT NULL REFERENCES groups(uuid) ON UPDATE CASCADE ON DELETE CASCADE,
    curator_uuid TEXT NOT NULL, --SAGA
    assigned_at  TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (group_uuid, curator_uuid)
);
CREATE INDEX idx_groups_curators_group_uuid ON "groups_curators" ("group_uuid");

CREATE TABLE IF NOT EXISTS "curators_hierarchy" (
    senior_uuid      TEXT NOT NULL,
    subordinate_uuid TEXT NOT NULL,
    created_at       TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (senior_uuid, subordinate_uuid),
    CHECK (senior_uuid != subordinate_uuid)
);
CREATE INDEX idx_curators_hierarchy_subordinate_uuid ON "curators_hierarchy" ("subordinate_uuid");