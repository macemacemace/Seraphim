package messages

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/macemacemace/seraphim/backend/internal/db"
)

var (
	ErrListingNotFound = errors.New("listing not found")
	ErrOwnListing      = errors.New("this listing is yours")
)

type messageStore interface {
	GetListing(ctx context.Context, id uuid.UUID) (db.GetListingRow, error)
	StartConversation(ctx context.Context, arg db.StartConversationParams) (db.Conversation, error)
}

type Service struct {
	store messageStore
}

func NewService(store messageStore) *Service {
	return &Service{store: store}
}

func (s *Service) Start(ctx context.Context, listingID, buyerID uuid.UUID) (db.Conversation, error) {
	listing, err := s.store.GetListing(ctx, listingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Conversation{}, ErrListingNotFound
		}
		return db.Conversation{}, fmt.Errorf("get listing: %w", err)
	}

	if listing.SellerID == buyerID {
		return db.Conversation{}, ErrOwnListing
	}

	conversation, err := s.store.StartConversation(ctx, db.StartConversationParams{
		ListingID: listingID,
		BuyerID:   buyerID,
		SellerID:  listing.SellerID,
	})

	if err != nil {
		return db.Conversation{}, fmt.Errorf("start conversation: %w", err)
	}

	return conversation, nil
}
