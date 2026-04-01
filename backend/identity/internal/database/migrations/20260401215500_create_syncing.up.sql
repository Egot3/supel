

CREATE TABLE IF NOT EXISTS sync_state (
    source_name TEXT PRIMARY KEY,
    last_synced_at TIMESTAMPZ NOT NULL,
    last_message_id UUID
);
