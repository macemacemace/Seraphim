-- +goose Up
CREATE TABLE conversations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id uuid NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    buyer_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seller_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_message_at timestamptz NOT NULL DEFAULT now(),
    buyer_last_read_at timestamptz,
    seller_last_read_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (listing_id, buyer_id),
    CHECK (buyer_id <> seller_id)
);
CREATE TABLE messages(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id uuid NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body text NOT NULL CHECK (
        char_length(body) BETWEEN 1 AND 2000
    ),
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX messages_conversation_created_idx ON messages (conversation_id, created_at);
-- +goose Down
DROP TABLE messages;
DROP TABLE conversations;
