-- Alter "news" table
ALTER TABLE news
ADD COLUMN deleted_at TIMESTAMPTZ;