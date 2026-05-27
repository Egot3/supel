CREATE TABLE IF NOT EXISTS periods (
    uuid       TEXT PRIMARY KEY DEFAULT (lower(
        hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-' || '4' || substr(hex(randomblob(2)),2) || '-' || 
        substr('89ab', (abs(random()) % 4) + 1, 1) || substr(hex(randomblob(2)),2) || '-' || 
        hex(randomblob(6))
    )),
    name       TEXT NOT NULL,
    position   INTEGER NOT NULL DEFAULT 0,
    start_time TEXT NOT NULL, 
    end_time   TEXT NOT NULL,
    deleted_at TEXT NULL,     
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS subjects (
    uuid         TEXT PRIMARY KEY DEFAULT (lower(
        hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-' || '4' || substr(hex(randomblob(2)),2) || '-' || 
        substr('89ab', (abs(random()) % 4) + 1, 1) || substr(hex(randomblob(2)),2) || '-' || 
        hex(randomblob(6))
    )),
    name         TEXT NOT NULL,
    deleted_at   TEXT NULL,
    created_at   TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at   TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE UNIQUE INDEX uq_subjects_name ON subjects (lower(name));

CREATE TABLE IF NOT EXISTS timetables (
    uuid       TEXT PRIMARY KEY DEFAULT (lower(
        hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-' || '4' || substr(hex(randomblob(2)),2) || '-' || 
        substr('89ab', (abs(random()) % 4) + 1, 1) || substr(hex(randomblob(2)),2) || '-' || 
        hex(randomblob(6))
    )),
    group_uuid TEXT NOT NULL,  -- saga 
    name       TEXT NOT NULL,
    assign_at  TEXT NULL,      
    revoke_at  TEXT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX idx_timetable_group_uuid ON timetables (group_uuid);

CREATE TABLE IF NOT EXISTS timetable_entries (
    uuid           TEXT PRIMARY KEY DEFAULT (lower(
        hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-' || '4' || substr(hex(randomblob(2)),2) || '-' || 
        substr('89ab', (abs(random()) % 4) + 1, 1) || substr(hex(randomblob(2)),2) || '-' || 
        hex(randomblob(6))
    )),
    timetable_uuid TEXT NOT NULL REFERENCES timetables(uuid) ON DELETE CASCADE,
    period_uuid    TEXT NOT NULL REFERENCES periods(uuid),
    day_of_week    INTEGER NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
    subject_uuid   TEXT NOT NULL REFERENCES subjects(uuid),
    place          TEXT NOT NULL,
    teacher_uuid   TEXT NULL,
    deleted_at     TEXT NULL,
    created_at     TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at     TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE UNIQUE INDEX uq_timetable_entries_active ON timetable_entries (timetable_uuid, period_uuid, day_of_week) WHERE deleted_at IS NULL;
CREATE INDEX idx_timetable_entries_timetable_uuid ON timetable_entries(timetable_uuid);

CREATE TABLE IF NOT EXISTS lessons (
    uuid                 TEXT PRIMARY KEY DEFAULT (lower(
        hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-' || '4' || substr(hex(randomblob(2)),2) || '-' || 
        substr('89ab', (abs(random()) % 4) + 1, 1) || substr(hex(randomblob(2)),2) || '-' || 
        hex(randomblob(6))
    )),
    timetable_entry_uuid TEXT NOT NULL REFERENCES timetable_entries(uuid) ON DELETE CASCADE,
    date                 TEXT NOT NULL,
    cancelled            INTEGER NOT NULL DEFAULT 0 CHECK (cancelled IN (0,1)),
    created_at           TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE (date, timetable_entry_uuid)
);

CREATE TABLE IF NOT EXISTS homework_attachments (
    uuid        TEXT PRIMARY KEY DEFAULT (lower(
        hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-' || '4' || substr(hex(randomblob(2)),2) || '-' || 
        substr('89ab', (abs(random()) % 4) + 1, 1) || substr(hex(randomblob(2)),2) || '-' || 
        hex(randomblob(6))
    )),
    lesson_uuid TEXT NOT NULL REFERENCES lessons(uuid) ON DELETE CASCADE,
    storage_key TEXT NOT NULL,
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);
