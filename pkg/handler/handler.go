package handler

import (
	"Todo_rest_api/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	api := router.Group("/api", h.userIdMiddleware)
	{
		lists := api.Group("/lists")
		{
			lists.GET("/:id", h.getListById)
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.DELETE("/:id", h.deleteList)
			lists.PUT("/:id", h.updateList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)

			}
		}

		items := api.Group("/items")
		{
			items.GET("/:id", h.getItemById)
			items.DELETE("/:id", h.deleteItem)
			items.PUT("/:id", h.updateItem)
		}

	}
	return router

}
