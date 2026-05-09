CREATE TYPE weekday AS ENUM (
    'Monday', 
    'Tuesday', 
    'Wednesday', 
    'Thursday', 
    'Friday', 
    'Saturday', 
    'Sunday'
);


-- Create "abstract_lessons" table
CREATE TABLE "abstract_lessons" (
  "uuid" uuid              PRIMARY KEY,
  "name" character varying NOT NULL,

  "deleted_at" TIMESTAMPTZ NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);


-- Create "concrete_lessons" table
CREATE TABLE "concrete_lessons" (
  "concrete_uuid" uuid NOT NULL,
  "abstract_uuid" uuid NOT NULL REFERENCES abstract_lessons(uuid) ON UPDATE CASCADE ON DELETE CASCADE,
  "teacher_uuid"  uuid NOT NULL,--SAGA
  "group_uuid"    uuid NOT NULL,--SAGA

  "homework_body_key" character varying GENERATED ALWAYS AS ('orgs/ETSEvilCorp/timetable/homework/body/' || "concrete_uuid"::text) STORED,

  "week_number"  SMALLINT NOT NULL,
  "day_of_week"  weekday NOT NULL
  "period"       SMALLINT NOT NULL,
  "year"         SMALLINT NOT NULL,

  "deleted_at" TIMESTAMPTZ NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY ("concrete_uuid")
);
CREATE INDEX idx_start_group_uuid   ON concrete_lessons ("week_number","group_uuid") WHERE deleted_at IS NULL;
CREATE INDEX idx_start_teacher_uuid ON concrete_lessons ("weel_number","teacher_uuid") WHERE deleted_at IS NULL;

-- Create "homework_attachment_keys"
CREATE TABLE "homework_attachments" (
  "file_uuid"     uuid PRIMARY KEY,
  "name"          character varying NOT NULL,
  "mime"          character varying NOT NULL,
  "concrete_uuid" uuid NOT NULL REFERENCES concrete_lessons(concrete_uuid) ON UPDATE NO ACTION ON DELETE CASCADE,

  "storage_key" character varying NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_concrete_uuid ON homework_attachments ("concrete_uuid");
