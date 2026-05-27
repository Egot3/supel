-- +migrate Up

ALTER TABLE IF EXISTS "periods"    ALTER COLUMN uuid DEFAULT gen_random_uuid();
ALTER TABLE IF EXISTS "subjects"   ALTER COLUMN uuid DEFAULT gen_random_uuid();
ALTER TABLE IF EXISTS "timetables" ALTER COLUMN uuid DEFAULT gen_random_uuid();
ALTER TABLE IF EXISTS "timetable_entries" ALTER COLUMN uuid DEFAULT gen_random_uuid();
ALTER TABLE IF EXISTS "lessons" ALTER COLUMN uuid DEFAULT gen_random_uuid();
ALTER TABLE IF EXISTS "homework_attachments" ALTER COLUMN uuid DEFAULT gen_random_uuid();