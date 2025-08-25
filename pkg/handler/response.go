package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Massage string `json:"massage"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

func NewErrorResponse(c *gin.Context, statusCode int, massage string) {
	logrus.Error(massage)
	c.AbortWithStatusJSON(statusCode, ErrorResponse{massage})
}
