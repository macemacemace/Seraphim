package listings

import (
	"errors"
	"log"
	"net/http"

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

type Handler struct {
	service *Service
}
