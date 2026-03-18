package http

import (
	"github.com/gin-gonic/gin"

	"tgPlanBot/internal/transport/http/handlers"
)

func registerRoutes(router *gin.Engine) {
	router.GET("/health", handlers.Health)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/tasks", handlers.GetTasks)
		}
	}
}
