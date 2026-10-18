package listings

import (
	"context"

	"github.com/google/uuid"
	"github.com/macemacemace/seraphim/backend/internal/db"
)

type listingStore interface {
	CreateListing(ctx context.Context, arg db.CreateListingParams) (db.CreateListingRow, error)
	GetListing(ctx context.Context, id uuid.UUID) (db.GetListingRow, error)
}

type Service struct {
	store listingStore
}

func NewService(store listingStore) *Service {
	return &Service{store: store}
}
