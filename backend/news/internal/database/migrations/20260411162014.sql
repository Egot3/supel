-- Create "news" table
CREATE TABLE "news" (
  "new_uuid" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_uuid" character varying NOT NULL,
  "caption" character varying NOT NULL,
  "body" character varying NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("new_uuid")
);
-- Create "news_images" table
CREATE TABLE "news_images" (
  "image_uuid" uuid NOT NULL DEFAULT gen_random_uuid(),
  "new_uuid" uuid NOT NULL,
  "file_key" character varying NOT NULL,
  "position" bigint NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("image_uuid"),
  CONSTRAINT "news_images_new_uuid_fkey" FOREIGN KEY ("new_uuid") REFERENCES "news" ("new_uuid") ON UPDATE NO ACTION ON DELETE NO ACTION
);
