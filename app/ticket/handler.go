package ticket

import (
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler(*UserService, *Storage) *Handler {
	return &Handler{}
}

func (handler *Handler) Booking(c *gin.Context) {}
