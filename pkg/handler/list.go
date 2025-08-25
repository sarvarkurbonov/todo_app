package handler

import (
	todo "Todo_rest_api"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetAllResponse struct {
	Data []todo.TodoList `json:"data"`
}

func (h *Handler) createList(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}

	var input todo.TodoList
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoList.CreateList(userId, input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, GetAllResponse{lists})
}
func (h *Handler) getListById(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	list, err := h.services.TodoList.GetById(userId, id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, list)
}
func (h *Handler) updateList(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	var input todo.UpdateList
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.TodoList.UpdateList(userId, id, input); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})

}
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	err = h.services.TodoList.DeleteList(userId, id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})

}
