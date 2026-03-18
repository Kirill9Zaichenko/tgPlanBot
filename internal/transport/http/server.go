package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"tgPlanBot/internal/config"
)

type Server struct {
	httpServer *http.Server
	host       string
	port       int
}

func NewServer(cfg *config.Config) *Server {
	router := NewRouter()

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		host: cfg.HTTP.Host,
		port: cfg.HTTP.Port,
	}
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

func NewRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	registerRoutes(router)

	return router
}
