CREATE TABLE IF NOT EXISTS news (
    new_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid NOT NULL,
    caption VARCHAR(255) NOT NULL,
    body TEXT DEFAULT NULL, --explicitly
    created_at TIMESTAMPZ NOT NULL
);

CREATE TABLE IF NOT EXISTS news_images {
    image_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    new_uuid UUID NOT NULL REFERENCES news(new_uuid) ON DELETE CASCADE,
    file_key TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPZ NOT NULL
};
