package listings

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/macemacemace/seraphim/backend/internal/db"
)

var (
	ErrInvalidType     = errors.New("type must be ITEM or SERVICE")
	ErrTitleRequired   = errors.New("title is required")
	ErrInvalidPrice    = errors.New("price must be 0 or more")
	ErrNotFound        = errors.New("listing not found")
	ErrInvalidLocation = errors.New("invalid location")
	ErrNotYourListing  = errors.New("you can only add photos to your own listings")
	ErrTooManyPhotos   = errors.New("a listing can have at most 5 photos")
)

type listingStore interface {
	CreateListing(ctx context.Context, arg db.CreateListingParams) (db.CreateListingRow, error)
	GetListing(ctx context.Context, id uuid.UUID) (db.GetListingRow, error)
	ListNearbyListings(ctx context.Context, arg db.ListNearbyListingsParams) ([]db.ListNearbyListingsRow, error)
	ListRecentListings(ctx context.Context) ([]db.ListRecentListingsRow, error)
	CreateListingPhoto(ctx context.Context, arg db.CreateListingPhotoParams) (db.ListingPhoto, error)
	ListListingPhotos(ctx context.Context, listingID uuid.UUID) ([]db.ListingPhoto, error)
}

type photoUploader interface {
	Upload(ctx context.Context, file io.Reader) (string, string, error)
}

type Service struct {
	store  listingStore
	photos photoUploader
}

func NewService(store listingStore, photos photoUploader) *Service {
	return &Service{store: store, photos: photos}

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

func (s *Service) Get(ctx context.Context, id uuid.UUID) (db.GetListingRow, error) {
	listing, err := s.store.GetListing(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.GetListingRow{}, ErrNotFound
		}

		return db.GetListingRow{}, fmt.Errorf("get listing: %w", err)
	}
	return listing, nil

}

func (s *Service) ListNearby(ctx context.Context, lat, lng float64) ([]db.ListNearbyListingsRow, error) {
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return nil, ErrInvalidLocation
	}

	listings, err := s.store.ListNearbyListings(ctx, db.ListNearbyListingsParams{
		Lat: lat,
		Lng: lng,
	})

	if err != nil {
		return nil, fmt.Errorf("list nearby listings: %w", err)
	}

	return listings, nil

}

func (s *Service) ListRecent(ctx context.Context) ([]db.ListRecentListingsRow, error) {

	listings, err := s.store.ListRecentListings(ctx)

	if err != nil {
		return nil, fmt.Errorf("list recent listings: %w", err)
	}

	return listings, nil

}

func (s *Service) AddPhoto(ctx context.Context, listingID, sellerID uuid.UUID, file io.Reader) (db.ListingPhoto, error) {
	listing, err := s.Get(ctx, listingID)
	if err != nil {
		return db.ListingPhoto{}, err
	}

	if listing.SellerID != sellerID {
		return db.ListingPhoto{}, ErrNotYourListing
	}

	existing, err := s.store.ListListingPhotos(ctx, listingID)
	if err != nil {
		return db.ListingPhoto{}, fmt.Errorf("list photos: %w", err)
	}

	if len(existing) >= 5 {
		return db.ListingPhoto{}, ErrTooManyPhotos
	}

	url, publicID, err := s.photos.Upload(ctx, file)
	if err != nil {
		return db.ListingPhoto{}, err
	}
	return s.store.CreateListingPhoto(ctx, db.CreateListingPhotoParams{
		ListingID: listingID,
		Url:       url,
		PublicID:  publicID,
		Position:  int32(len(existing)),
	})

}

func (s *Service) ListingPhotos(ctx context.Context, listingID uuid.UUID) ([]db.ListingPhoto, error) {

	photos, err := s.store.ListListingPhotos(ctx, listingID)

	if err != nil {
		return nil, fmt.Errorf("list photos: %w", err)
	}

	return photos, nil
}
