-- +migrate Up

ALTER IF EXISTS TABLE "periods"    ALTER COLUMN uuid DEFAULT gen_random_uuid();
ALTER IF EXISTS TABLE "subjects"   ALTER COLUMN uuid DEFAULT gen_random_uuid();
ALTER IF EXISTS TABLE "timetables" ALTER COLUMN uuid DEFAULT gen_random_uuid();
ALTER IF EXISTS TABLE "timetable_entries" ALTER COLUMN uuid DEFAULT gen_random_uuid();
ALTER IF EXISTS TABLE "lessons" ALTER COLUMN uuid DEFAULT gen_random_uuid();
ALTER IF EXISTS TABLE "homework_attachments" ALTER COLUMN uuid DEFAULT gen_random_uuid();