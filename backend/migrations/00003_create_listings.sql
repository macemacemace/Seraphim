-- +goose Up
CREATE TABLE listings (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type text NOT NULL CHECK (type IN ('ITEM', 'SERVICE')),
    title text NOT NULL,
    description text NOT NULL DEFAULT '',
    price integer NOT NULL CHECK (price >= 0),
    status text NOT NULL CHECK (status IN('ACTIVE', 'SOLD')) DEFAULT 'ACTIVE',
    location geography(Point, 4326) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX listings_location_idx ON listings USING GIST (location);
CREATE INDEX listings_seller_id_idx ON listings(seller_id);
-- +goose Down
DROP TABLE listings;
