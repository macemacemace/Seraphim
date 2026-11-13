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
    updated_at,
    COALESCE(
        (
            SELECT url
            FROM listing_photos
            WHERE listing_photos.listing_id = listings.id
            ORDER BY position
            LIMIT 1
        ), ''
    )::text AS cover_url
FROM listings
WHERE ST_DWithin(
        location,
        ST_SetSRID(ST_MakePoint(@lng::float8, @lat::float8), 4326)::geography,
        50000
    )
    AND status = 'ACTIVE'
    AND (
        @search::text = ''
        OR title ILIKE '%' || @search || '%'
        OR description ILIKE '%' || @search || '%'
    )
    AND (
        @listing_type::text = ''
        OR type = @listing_type
    )
    AND price >= @min_price::int
    AND (
        @max_price::int = -1
        OR price <= @max_price
    )
ORDER BY CASE
        WHEN @sort::text = 'price_asc' THEN price
    END ASC,
    CASE
        WHEN @sort::text = 'price_desc' THEN price
    END DESC,
    CASE
        WHEN @sort::text = 'oldest' THEN listings.created_at
    END ASC,
    listings.created_at DESC
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
    updated_at,
    COALESCE(
        (
            SELECT url
            FROM listing_photos
            WHERE listing_photos.listing_id = listings.id
            ORDER BY position
            LIMIT 1
        ), ''
    )::text AS cover_url
FROM listings
WHERE status = 'ACTIVE'
    AND (
        @search::text = ''
        OR title ILIKE '%' || @search || '%'
        OR description ILIKE '%' || @search || '%'
    )
    AND (
        @listing_type::text = ''
        OR type = @listing_type
    )
    AND price >= @min_price::int
    AND (
        @max_price::int = -1
        OR price <= @max_price
    )
ORDER BY CASE
        WHEN @sort::text = 'price_asc' THEN price
    END ASC,
    CASE
        WHEN @sort::text = 'price_desc' THEN price
    END DESC,
    CASE
        WHEN @sort::text = 'oldest' THEN listings.created_at
    END ASC,
    listings.created_at DESC
LIMIT 50;
-- name: ListListingsBySeller :many
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
    updated_at,
    COALESCE (
        (
            SELECT url
            FROM listing_photos
            WHERE listing_photos.listing_id = listing_id
            ORDER BY position
            LIMIT 1
        ), ''
    )::text AS cover_url
FROM listings
WHERE seller_id = @seller_id
ORDER BY created_at DESC;
