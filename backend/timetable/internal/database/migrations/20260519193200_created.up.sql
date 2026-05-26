-- +migrate Up
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS "periods" (
    uuid       UUID PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    position   SMALLINT NOT NULL DEFAULT 0,
    start_time TIME NOT NULL,
    end_time   TIME NOT NULL,

    "deleted_at" TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT start_end_no_overlap CHECK (start_time < end_time)
);

CREATE TABLE IF NOT EXISTS "subjects" (
    "uuid" UUID PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,

    "deleted_at" TIMESTAMPTZ NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX uq_subjects_name ON subjects (lower("name"));

CREATE TABLE IF NOT EXISTS "timetables" (
    uuid UUID PRIMARY KEY,
    group_uuid UUID NOT NULL, --SAGA
    name       VARCHAR(65) NOT NULL,
    assign_at DATE,
    revoke_at  DATE,

    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),    
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT timetable_valid_range CHECK (
        assign_at IS NULL OR revoke_at IS NULL OR assign_at<=revoke_at
    ),
    CONSTRAINT no_overlappingGroupTimetables EXCLUDE USING gist (
        group_uuid WITH =,
        daterange(assign_at, revoke_at, '[]') WITH &&
    )
);
CREATE INDEX idx_timetable_group_uuid ON "timetables" ("group_uuid");

CREATE TABLE IF NOT EXISTS timetable_entries (
    uuid           UUID PRIMARY KEY,
    timetable_uuid UUID NOT NULL REFERENCES timetables(uuid) ON DELETE CASCADE,
    period_uuid    UUID NOT NULL REFERENCES periods(uuid),
    day_of_week    SMALLINT NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
    subject_uuid   UUID NOT NULL REFERENCES subjects(uuid),
    place          VARCHAR(65) NOT NULL,
    teacher_uuid   UUID NULL,

    "deleted_at" TIMESTAMPTZ NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX uq_timetable_entries_active ON timetable_entries (timetable_uuid, period_uuid, day_of_week) WHERE deleted_at IS NULL;
CREATE INDEX idx_timetable_entries_timetable_uuid ON timetable_entries(timetable_uuid);

CREATE TABLE IF NOT EXISTS lessons (
    uuid UUID PRIMARY KEY,
    timetable_entry_uuid UUID NOT NULL REFERENCES timetable_entries(uuid) ON DELETE CASCADE,
    date      DATE NOT NULL,
    cancelled BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (date, timetable_entry_uuid)
);

CREATE TABLE IF NOT EXISTS "homework_attachments" (
  "uuid"        uuid NOT NULL PRIMARY KEY,
  "lesson_uuid" uuid NOT NULL REFERENCES lessons(uuid) ON DELETE CASCADE,

  "storage_key" character varying NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);