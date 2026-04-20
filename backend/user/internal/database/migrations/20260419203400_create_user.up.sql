-- Create "news" table
CREATE TABLE "news" (
  "uuid" uuid NOT NULL,
  "nickname" character varying NOT NULL,
  "avatar_key" uuid NULL,
  "description" character varying NULL,

  "status" character varying NULL,
  "status_expiration" TIMESTAMPTZ NULL,
  "status_reaction_key" TIMESTAMPTZ NULL,

  "deleted_at" TIMESTAMPTZ NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("uuid")
);
