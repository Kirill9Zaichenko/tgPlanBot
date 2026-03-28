package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	moderationapp "tgPlanBot/internal/app/moderation"
	taskapp "tgPlanBot/internal/app/task"
	"tgPlanBot/internal/config"
)

type Server struct {
	httpServer        *http.Server
	host              string
	port              int
	taskService       *taskapp.Service
	moderationService *moderationapp.Service
}

func NewServer(
	cfg *config.Config,
	taskService *taskapp.Service,
	moderationService *moderationapp.Service,
) *Server {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	s := &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
			Handler: router,
		},
		host:              cfg.HTTP.Host,
		port:              cfg.HTTP.Port,
		taskService:       taskService,
		moderationService: moderationService,
	}

	s.registerRoutes(router)

	return s
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Address() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}
