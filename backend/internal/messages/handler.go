package messages

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type startRequest struct {
	ListingID string `json:"listingId"`
}

func (h *Handler) Start(c *gin.Context) {
	buyerID := c.MustGet("userID").(uuid.UUID)

	var req startRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	listingID, err := uuid.Parse(req.ListingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid listing id"})
		return
	}
	conversation, err := h.service.Start(c.Request.Context(), listingID, buyerID)

	if err != nil {
		if errors.Is(err, ErrListingNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, ErrOwnListing) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Println("start conversation: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": conversation.ID})

}
