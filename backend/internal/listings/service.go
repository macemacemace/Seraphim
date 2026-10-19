package listings

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/macemacemace/seraphim/backend/internal/db"
)

var (
	ErrInvalidType   = errors.New("type must be ITEM or SERVICE")
	ErrTitleRequired = errors.New("title is required")
	ErrInvalidPrice  = errors.New("price must be 0 or more")
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

type CreateInput struct {
	SellerID    uuid.UUID
	Type        string
	Title       string
	Description string
	Price       int32
	Lat         float64
	Lng         float64
}

func (s *Service) Create(ctx context.Context, input CreateInput) (db.CreateListingRow, error) {

	input.Title = strings.TrimSpace(input.Title)

	if input.Title == "" {
		return db.CreateListingRow{}, ErrTitleRequired
	}

	if input.Price < 0 {
		return db.CreateListingRow{}, ErrInvalidPrice
	}

	if input.Type != "ITEM" && input.Type != "SERVICE" {
		return db.CreateListingRow{}, ErrInvalidType
	}

	lat := math.Round(input.Lat*100) / 100
	lng := math.Round(input.Lng*100) / 100

	listing, err := s.store.CreateListing(ctx, db.CreateListingParams{

		SellerID:    input.SellerID,
		Type:        input.Type,
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Lat:         lat,
		Lng:         lng,
	})
	if err != nil {
		return db.CreateListingRow{}, fmt.Errorf("create listing: %w", err)

	}

	return listing, nil

}
