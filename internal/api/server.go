package api

import (
	"context"
	"log"
	"net/http"

	"github.com/andrearcaina/whisp/internal/config"
	"github.com/andrearcaina/whisp/internal/db"
	http2 "github.com/andrearcaina/whisp/internal/handlers/http"
	"github.com/andrearcaina/whisp/internal/handlers/ws"
	"github.com/gin-gonic/gin"
)

type Server struct {
	HTTP *http.Server
	DB   *db.Database
}

func NewWebServer() *Server {
	cfg := config.NewConfig()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// init the database connection
	database, err := db.New(cfg)
	if err != nil {
		return nil
	}

	// init the WebSocket hub
	hub := ws.NewHub()
	go hub.Run()

	// now create an http handler with the hub and db
	h := http2.NewHandler(database, hub, cfg.TenorAPIKey)

	// create the http server
	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: h.NewRouter(), // set the router from the handler
	}

	return &Server{
		HTTP: srv,
		DB:   database,
	}
}

func (s *Server) Run() error {
	log.Printf("Starting whisp app on %s\n", s.HTTP.Addr)
	return s.HTTP.ListenAndServe()
}

func (s *Server) GracefulShutdown(ctx context.Context) error {
	log.Println("Shutting down whisp gracefully...")

	if s.DB != nil {
		log.Println("Database connection closed.")
		s.DB.Close()
	}

	err := s.HTTP.Shutdown(ctx)
	if err != nil {
		return err
	}

	log.Println("whisp shutdown complete.")
	return nil
}
