-- name: CreateListing :one
INSERT INTO listings (
        seller_id,
        type,
        title,
        description,
        price,
        location
    )
VALUES (
        @seller_id,
        @type,
        @title,
        @description,
        @price,
        ST_SetSRID(ST_MakePoint(@lng::float8, @lat::float8), 4326)::geography
    )
RETURNING id,
    seller_id,
    type,
    title,
    description,
    price,
    status,
    ST_Y(location::geometry)::float8 AS lat,
    ST_X(location::geometry)::float8 AS lng,
    created_at,
    updated_at;
-- name: GetListing :one
SELECT id,
    seller_id,
    type,
    title,
    description,
    price,
    status,
    ST_Y(location::geometry)::float8 AS lat,
    ST_X(location::geometry)::float8 AS lng,
    created_at,
    updated_at
FROM listings
WHERE id = @id;
-- name: ListNearbyListings :many
SELECT id,
    seller_id,
    type,
    title,
    description,
    price,
    status,
    ST_Y(location::geometry)::float8 AS lat,
    ST_X(location::geometry)::float8 AS lng,
    created_at,
    updated_at
FROM listings
WHERE ST_DWithin(
        location,
        ST_SetSRID(ST_MakePoint(@lng::float8, @lat::float8), 4326)::geography,
        50000
    )
    AND status = 'ACTIVE'
ORDER BY created_at DESC
LIMIT 50;
-- name: ListRecentListings :many
SELECT id,
    seller_id,
    type,
    title,
    description,
    price,
    status,
    ST_Y(location::geometry)::float8 AS lat,
    ST_X(location::geometry)::float8 AS lng,
    created_at,
    updated_at
FROM listings
WHERE status = 'ACTIVE'
ORDER BY created_at DESC
LIMIT 50;
