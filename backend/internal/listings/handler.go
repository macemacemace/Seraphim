package listings

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"io"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createListingRequest struct {
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       int32   `json:"price"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
}

type createListingResponse struct {
	ID          uuid.UUID `json:"id"`
	SellerID    uuid.UUID `json:"sellerId"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int32     `json:"price"`
	Lat         float64   `json:"lat"`
	Lng         float64   `json:"lng"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	sellerID := c.MustGet("userID").(uuid.UUID)
	var req createListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	listing, err := h.service.Create(c.Request.Context(), CreateInput{

		SellerID:    sellerID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Lat:         req.Lat,
		Lng:         req.Lng,
	})

	if err != nil {
		if errors.Is(err, ErrInvalidType) || errors.Is(err, ErrInvalidPrice) || errors.Is(err, ErrTitleRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Println("create listing:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, createListingResponse{
		ID:          listing.ID,
		SellerID:    listing.SellerID,
		Type:        listing.Type,
		Title:       listing.Title,
		Description: listing.Description,
		Price:       listing.Price,
		Lat:         listing.Lat,
		Lng:         listing.Lng,
	})

}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid listing id"})
		return
	}

	listing, err := h.service.Get(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		log.Println("get listing:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return

	}

	c.JSON(http.StatusOK, createListingResponse{
		ID:          listing.ID,
		SellerID:    listing.SellerID,
		Type:        listing.Type,
		Title:       listing.Title,
		Description: listing.Description,
		Price:       listing.Price,
		Lat:         listing.Lat,
		Lng:         listing.Lng,
	})
}

func (h *Handler) ListNearby(c *gin.Context) {

	if c.Query("lat") == "" && c.Query("lng") == "" {
		listings, err := h.service.ListRecent(c.Request.Context())

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		response := make([]createListingResponse, 0, len(listings))
		for _, l := range listings {
			response = append(response, createListingResponse{
				ID:          l.ID,
				SellerID:    l.SellerID,
				Type:        l.Type,
				Title:       l.Title,
				Description: l.Description,
				Price:       l.Price,
				Lat:         l.Lat,
				Lng:         l.Lng,
			})

		}
		log.Println("list recent listings:", err)
		c.JSON(http.StatusOK, response)
		return
	}
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lng are required numbers"})
		return
	}
	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lng are required numbers"})
		return
	}

	listings, err := h.service.ListNearby(c.Request.Context(), lat, lng)
	if err != nil {
		if errors.Is(err, ErrInvalidLocation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Println("list nearby listings:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	response := make([]createListingResponse, 0, len(listings))
	for _, l := range listings {
		response = append(response, createListingResponse{
			ID:          l.ID,
			SellerID:    l.SellerID,
			Type:        l.Type,
			Title:       l.Title,
			Description: l.Description,
			Price:       l.Price,
			Lat:         l.Lat,
			Lng:         l.Lng,
		})
	}
	c.JSON(http.StatusOK, response)

}

type Handler struct {
	service *Service
}

const maxPhotoSize = 5 << 20

func (h *Handler) AddPhoto(c *gin.Context) {
	sellerID := c.MustGet("userID").(uuid.UUID)

	listingID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid listing id"})
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxPhotoSize+1024)

	fileHeader, err := c.FormFile("photo")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo is required and must be under 5MB"})
		return
	}

	file, err := fileHeader.Open()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read photo"})
		return
	}

	defer file.Close()

	head := make([]byte, 512)
	n, err := file.Read(head)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read photo"})
		return
	}

	contentType := http.DetectContentType(head[:n])
	if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo must be JPG, PNG or WebP"})
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	photo, err := h.service.AddPhoto(c.Request.Context(), listingID, sellerID, file)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrNotYourListing) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		if errors.Is(err, ErrTooManyPhotos) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		log.Println("add photo:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, photoResponse{
		ID:       photo.ID,
		URL:      photo.Url,
		Position: photo.Position,
	})

}
