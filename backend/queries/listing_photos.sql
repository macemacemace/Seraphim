-- name: CreateListingPhoto :one
INSERT INTO listing_photos (listing_id, url, public_id, position)
VALUES (@listing_id, @url, @public_id, @position)
RETURNING *;
-- name: ListListingPhotos :many
SELECT *
FROM listing_photos
WHERE listing_id = @listing_id
ORDER BY position;
