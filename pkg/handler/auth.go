package handler

import (
	todo "Todo_rest_api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})

}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {

	var input signInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	token, err := h.services.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"token": token})

}
