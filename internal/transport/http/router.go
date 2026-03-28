package http

import (
	"tgPlanBot/internal/transport/http/handlers"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerRoutes(router *gin.Engine) {
	router.GET("/health", handlers.Health)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/tasks", handlers.GetTasks(s.taskService))
			v1.POST("/tasks", handlers.CreateTask(s.taskService))

			v1.GET("/inbox", handlers.GetInbox(s.moderationService))
			v1.POST("/tasks/:id/accept", handlers.AcceptTask(s.moderationService))
			v1.POST("/tasks/:id/reject", handlers.RejectTask(s.moderationService))
		}
	}
}
