package handler

import (
	todo "Todo_rest_api"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createItem(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}
	var input todo.TodoItem
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.CreateItem(userId, listId, input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})

}
func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}
	items, err := h.services.TodoItem.GetAll(userId, listId)

	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)

}
func (h *Handler) getItemById(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}
	item, err := h.services.TodoItem.GetById(c.Request.Context(), userId, itemId)

	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *Handler) updateItem(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	var input todo.UpdateItem
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.services.TodoItem.UpdateItem(c.Request.Context(), userId, itemId, input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})

}

func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}
	err = h.services.TodoItem.Delete(c.Request.Context(), userId, itemId)

	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})

}
