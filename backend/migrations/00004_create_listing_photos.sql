-- +goose Up
CREATE TABLE listing_photos (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id uuid NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    url text NOT NULL,
    public_id text NOT NULL,
    position integer NOT NULL CHECK (
        position BETWEEN 0 AND 4
    ),
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (listing_id, position)
);
-- +goose Down
DROP TABLE listing_photos;
