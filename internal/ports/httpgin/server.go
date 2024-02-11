package httpgin

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"ads-server/internal/app"
)

type Server struct {
	port string
	app  *gin.Engine
}

// loggerMW represents simple logger handler with useful info about request and time to request
func loggerMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		status := c.Writer.Status()
		log.Println("\n[\n\r", "Time:", latency, "\nMethod used:",
			c.Request.Method, "\nRequest to:", c.Request.URL.Path, "\nStatus code:", status, "\n]")
	}
}

func NewHTTPServer(port string, a app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	api := router.Group("api/v1")
	s := &http.Server{Addr: port, Handler: router}
	//api := s.Handler.Group("/api/v1")
	api.Use(loggerMW(), gin.Recovery())
	AppRouter(api, a)
	return s

}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}

// Run returns function to start HTTP server on a port given and implements graceful shutdown principle
func Run(ctx context.Context, a app.AdRepository, u app.UserRepository, httpPort string) func() error {
	return func() error {
		httpServer := NewHTTPServer(httpPort, app.NewApp(a, u))

		errCh := make(chan error)

		defer func() {
			shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(shCtx); err != nil {
				log.Printf("can't close http server listening on %s: %s", httpServer.Addr, err.Error())
			}

			close(errCh)
		}()

		go func() {
			if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("http server can't listen and serve requests: %w", err)
		}
	}
}
