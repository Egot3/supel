-- +migrate Up

-- periods
CREATE TABLE IF NOT EXISTS periods (
    uuid       TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    position   INTEGER NOT NULL DEFAULT 0,
    start_time TEXT NOT NULL,
    end_time   TEXT NOT NULL,

    deleted_at TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (start_time < end_time)
);

-- subjects
CREATE TABLE IF NOT EXISTS subjects (
    uuid TEXT PRIMARY KEY,
    name TEXT NOT NULL,

    deleted_at TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX uq_subjects_name ON subjects (lower(name));

-- timetables
CREATE TABLE IF NOT EXISTS timetables (
    uuid       TEXT PRIMARY KEY,
    group_uuid TEXT NOT NULL,   -- SAGA
    name       TEXT NOT NULL,
    assign_at  TEXT,
    revoke_at  TEXT,

    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (assign_at IS NULL OR revoke_at IS NULL OR assign_at <= revoke_at)
);
CREATE INDEX idx_timetable_group_uuid ON timetables (group_uuid);

-- timetable_entries
CREATE TABLE IF NOT EXISTS timetable_entries (
    uuid           TEXT PRIMARY KEY,
    timetable_uuid TEXT NOT NULL REFERENCES timetables(uuid) ON DELETE CASCADE,
    period_uuid    TEXT NOT NULL REFERENCES periods(uuid),
    day_of_week    INTEGER NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
    subject_uuid   TEXT NOT NULL REFERENCES subjects(uuid),
    place          TEXT NOT NULL,
    teacher_uuid   TEXT,

    deleted_at TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX uq_timetable_entries_active
    ON timetable_entries (timetable_uuid, period_uuid, day_of_week)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_timetable_entries_timetable_uuid ON timetable_entries(timetable_uuid);

-- lessons
CREATE TABLE IF NOT EXISTS lessons (
    uuid                  TEXT PRIMARY KEY,
    timetable_entry_uuid  TEXT NOT NULL REFERENCES timetable_entries(uuid) ON DELETE CASCADE,
    date                  TEXT NOT NULL,
    cancelled             INTEGER NOT NULL DEFAULT 0 CHECK (cancelled IN (0, 1)),
    created_at            TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (date, timetable_entry_uuid)
);

-- homework_attachments
CREATE TABLE IF NOT EXISTS homework_attachments (
    name         TEXT NOT NULL,
    mime         TEXT NOT NULL,
    lesson_uuid  TEXT NOT NULL REFERENCES lessons(uuid) ON DELETE CASCADE,

    storage_key  TEXT NOT NULL,
    created_at   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (lesson_uuid, name)
);