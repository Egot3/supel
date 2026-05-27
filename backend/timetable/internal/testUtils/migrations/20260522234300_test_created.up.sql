CREATE TABLE IF NOT EXISTS periods (
    uuid       TEXT PRIMARY KEY NOT NULL,
    name       TEXT NOT NULL,
    position   INTEGER NOT NULL DEFAULT 0,
    start_time TEXT NOT NULL, 
    end_time   TEXT NOT NULL,
    deleted_at TEXT NULL,     
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    CHECK (start_time < end_time)
);

CREATE TABLE IF NOT EXISTS subjects (
    uuid         TEXT PRIMARY KEY NOT NULL,
    name         TEXT NOT NULL,
    deleted_at   TEXT NULL,
    created_at   TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at   TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE UNIQUE INDEX uq_subjects_name ON subjects (lower(name));

CREATE TABLE IF NOT EXISTS timetables (
    uuid       TEXT PRIMARY KEY NOT NULL,
    group_uuid TEXT NOT NULL,  -- saga 
    name       TEXT NOT NULL,
    assign_at  TEXT NULL,      
    revoke_at  TEXT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    CHECK (
        assign_at IS NULL OR revoke_at IS NULL OR assign_at <= revoke_at
    )
);
CREATE INDEX idx_timetable_group_uuid ON timetables (group_uuid);

CREATE TABLE IF NOT EXISTS timetable_entries (
    uuid           TEXT PRIMARY KEY NOT NULL,
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
    uuid                 TEXT PRIMARY KEY NOT NULL,
    timetable_entry_uuid TEXT NOT NULL REFERENCES timetable_entries(uuid) ON DELETE CASCADE,
    date                 TEXT NOT NULL,
    cancelled            INTEGER NOT NULL DEFAULT 0 CHECK (cancelled IN (0,1)),
    created_at           TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE (date, timetable_entry_uuid)
);

CREATE TABLE IF NOT EXISTS homework_attachments (
    uuid        TEXT PRIMARY KEY NOT NULL,
    lesson_uuid TEXT NOT NULL REFERENCES lessons(uuid) ON DELETE CASCADE,
    storage_key TEXT NOT NULL,
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TRIGGER IF NOT EXISTS trg_periods_uuid
    BEFORE INSERT ON periods
    FOR EACH ROW
    WHEN NEW.uuid IS NULL
BEGIN
    SELECT
        lower(hex(randomblob(4))) || '-' ||
        lower(hex(randomblob(2))) || '-' ||
        '4' || substr(lower(hex(randomblob(2))),2) || '-' ||
        substr('89ab', (abs(random()) % 4) + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' ||
        lower(hex(randomblob(6)))
    INTO NEW.uuid;
END;

CREATE TRIGGER IF NOT EXISTS trg_subjects_uuid
    BEFORE INSERT ON subjects
    FOR EACH ROW WHEN NEW.uuid IS NULL
BEGIN
    SELECT lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-' || '4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab', (abs(random()) % 4) + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6))) INTO NEW.uuid;
END;

CREATE TRIGGER IF NOT EXISTS trg_timetables_uuid
    BEFORE INSERT ON timetables
    FOR EACH ROW WHEN NEW.uuid IS NULL
BEGIN
    SELECT lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-' || '4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab', (abs(random()) % 4) + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6))) INTO NEW.uuid;
END;

CREATE TRIGGER IF NOT EXISTS trg_timetable_entries_uuid
    BEFORE INSERT ON timetable_entries
    FOR EACH ROW WHEN NEW.uuid IS NULL
BEGIN
    SELECT lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-' || '4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab', (abs(random()) % 4) + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6))) INTO NEW.uuid;
END;

CREATE TRIGGER IF NOT EXISTS trg_lessons_uuid
    BEFORE INSERT ON lessons
    FOR EACH ROW WHEN NEW.uuid IS NULL
BEGIN
    SELECT lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-' || '4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab', (abs(random()) % 4) + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6))) INTO NEW.uuid;
END;

CREATE TRIGGER IF NOT EXISTS trg_homework_attachments_uuid
    BEFORE INSERT ON homework_attachments
    FOR EACH ROW WHEN NEW.uuid IS NULL
BEGIN
    SELECT lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-' || '4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab', (abs(random()) % 4) + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6))) INTO NEW.uuid;
END;

CREATE TRIGGER IF NOT EXISTS trg_timetables_no_overlap_insert
    BEFORE INSERT ON timetables
    WHEN NEW.assign_at IS NOT NULL AND NEW.revoke_at IS NOT NULL
BEGIN
    SELECT RAISE(ABORT, 'Overlapping timetable range for the same group')
    WHERE EXISTS (
        SELECT 1 FROM timetables
        WHERE group_uuid = NEW.group_uuid
          AND assign_at IS NOT NULL
          AND revoke_at IS NOT NULL
          AND assign_at <= NEW.revoke_at
          AND revoke_at >= NEW.assign_at
    );
END;

CREATE TRIGGER IF NOT EXISTS trg_timetables_no_overlap_update
    BEFORE UPDATE ON timetables
    WHEN NEW.assign_at IS NOT NULL AND NEW.revoke_at IS NOT NULL
BEGIN
    SELECT RAISE(ABORT, 'Overlapping timetable range for the same group')
    WHERE EXISTS (
        SELECT 1 FROM timetables
        WHERE group_uuid = NEW.group_uuid
          AND uuid != OLD.uuid
          AND assign_at IS NOT NULL
          AND revoke_at IS NOT NULL
          AND assign_at <= NEW.revoke_at
          AND revoke_at >= NEW.assign_at
    );
END;