package listings

import (
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

	_ = listing
	_ = err

}

type Handler struct {
	service *Service
}
