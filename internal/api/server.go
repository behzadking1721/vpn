package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"vpnclient/internal/managers"
)

// Server represents the API server
type Server struct {
	router          *mux.Router
	serverManager   *managers.ServerManager
	connectionManager *managers.ConnectionManager
	addr            string
	httpServer      *http.Server
}

// NewServer creates a new API server
func NewServer(
	addr string,
	serverManager *managers.ServerManager,
	connectionManager *managers.ConnectionManager,
) *Server {
	s := &Server{
		router:            mux.NewRouter(),
		serverManager:     serverManager,
		connectionManager: connectionManager,
		addr:              addr,
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	api := s.router.PathPrefix("/api").Subrouter()

	// Server management endpoints
	api.HandleFunc("/servers", s.listServers).Methods("GET")
	api.HandleFunc("/servers/enabled", s.listEnabledServers).Methods("GET")
	api.HandleFunc("/servers", s.addServer).Methods("POST")
	api.HandleFunc("/servers/{id}", s.getServer).Methods("GET")
	api.HandleFunc("/servers/{id}", s.updateServer).Methods("PUT")
	api.HandleFunc("/servers/{id}", s.deleteServer).Methods("DELETE")
	api.HandleFunc("/servers/{id}/enable", s.enableServer).Methods("POST")
	api.HandleFunc("/servers/{id}/disable", s.disableServer).Methods("POST")
	api.HandleFunc("/servers/{id}/ping", s.updatePing).Methods("PUT")
	api.HandleFunc("/servers/{id}/test-ping", s.testServerPing).Methods("POST")
	api.HandleFunc("/servers/test-all-ping", s.testAllServersPing).Methods("POST")
	api.HandleFunc("/servers/best", s.getBestServer).Methods("GET")

	// Subscription management endpoints
	api.HandleFunc("/subscriptions", s.getAllSubscriptions).Methods("GET")
	api.HandleFunc("/subscriptions", s.addSubscription).Methods("POST")
	api.HandleFunc("/subscriptions/{id}", s.getSubscription).Methods("GET")
	api.HandleFunc("/subscriptions/{id}", s.updateSubscription).Methods("PUT")
	api.HandleFunc("/subscriptions/{id}", s.deleteSubscription).Methods("DELETE")
	api.HandleFunc("/subscriptions/{id}/update", s.updateSubscriptionServers).Methods("POST")

	// Connection management endpoints
	api.HandleFunc("/connect", s.connect).Methods("POST")
	api.HandleFunc("/connect/fastest", s.connectFastest).Methods("POST")
	api.HandleFunc("/connect/best", s.connectBest).Methods("POST")
	api.HandleFunc("/disconnect", s.disconnect).Methods("POST")
	api.HandleFunc("/status", s.getStatus).Methods("GET")
	api.HandleFunc("/stats", s.getStats).Methods("GET")

	// Health check
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")

	// CORS middleware
	s.router.Use(corsMiddleware)
}

// Start starts the API server
func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Addr:         s.addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("API Server starting on %s\n", s.addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	fmt.Println("Server exited")
	return nil
}

// GetRouter returns the router for testing purposes
func (s *Server) GetRouter() *mux.Router {
	return s.router
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// healthCheck returns server health status
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}