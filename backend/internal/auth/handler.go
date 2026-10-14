package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.service.Register(c.Request.Context(), req.Email, req.Password, req.Name)

	if err != nil {
		if errors.Is(err, ErrInvalidEmail) || errors.Is(err, ErrWeakPassword) || errors.Is(err, ErrNameRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrEmailTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

			return
		}
		log.Printf("register: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return

	}
	c.JSON(http.StatusCreated, registerResponse{ID: user.ID, Email: user.Email, Name: user.Name})

}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})

		return
	}
	token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return

		}
		log.Printf("login: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return

	}
	c.JSON(http.StatusOK, loginResponse{Token: token})

}

func (h *Handler) Me(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	c.JSON(http.StatusOK, gin.H{"userID": userID})
}
