package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) userIdMiddleware(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	userId, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	c.Set("userId", userId)
	c.Next()

}

func GetUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("userId")
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "User id not found")
		return 0, errors.New("user id not found")

	}
	IdInt, ok := id.(int)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "User id not found")
		return 0, errors.New("user id not found")
	}
	return IdInt, nil
}
