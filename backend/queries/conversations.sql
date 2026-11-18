-- name: StartConversation :one
INSERT INTO conversations (listing_id, buyer_id, seller_id)
VALUES (@listing_id, @buyer_id, @seller_id) ON CONFLICT (listing_id, buyer_id) DO
UPDATE
SET listing_id = EXCLUDED.listing_id
RETURNING *;
-- name: GetConversation :one
SELECT *
FROM conversations
WHERE id = @id;
-- name: CreateMessage :one
INSERT INTO messages (conversation_id, sender_id, body)
VALUES (@conversation_id, @sender_id, @body)
RETURNING *;
-- name: TouchConversation :exec
UPDATE conversations
SET last_message_at = now()
WHERE id = @id;
-- name: ListMessages :many
SELECT *
FROM messages
WHERE conversation_id = @conversation_id
ORDER BY created_at ASC;
