package server

import (
	"context"
	config "github.com/romangurevitch/redis-cache-go"
	"github.com/romangurevitch/redis-cache-go/cache"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"time"
)

type server struct {
	cache  cache.Cache
	router *http.ServeMux
	logger *log.Logger
	apiUrl *url.URL
}

// Create new cached contact server
func NewContactServer(apiPath string, cache cache.Cache) (*server, error) {
	apiUrl, err := url.Parse(apiPath)
	if err != nil {
		return nil, err
	}
	return &server{
		cache:  cache,
		router: http.NewServeMux(),
		logger: log.New(os.Stdout, "http: ", log.LstdFlags),
		apiUrl: apiUrl,
	}, nil
}

// Start the server
func (s *server) Start() {
	s.logger.Println("Caching contact server is starting...")

	s.routes()
	server := &http.Server{
		Addr:         ":" + config.HttpPort,
		Handler:      s.router,
		ErrorLog:     s.logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		s.logger.Println("Caching contact server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			s.logger.Fatalf("Could not gracefully shutdown the contact: %v\n", err)
		}
		close(done)
	}()

	s.logger.Println("Caching contact server is ready to handle requests at localhost:" + config.HttpPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatalf("Could not listen on %s: %v\n", "localhost:"+config.HttpPort, err)
	}

	<-done
	s.logger.Println("Server stopped")
}

// Redirect request to the provided url
func (s *server) redirect(url *url.URL, w http.ResponseWriter, r *http.Request, interceptResponse func(*http.Response) error) {
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = interceptResponse
	proxy.ServeHTTP(w, r)
}
