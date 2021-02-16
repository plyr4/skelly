package router

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/tomb.v2"
)

// Run executes router to serve http for the application
func Run(port string) error {

	// router configurations
	router := gin.New()
	router.Use(gin.Recovery())

	// health endpoint
	router.GET("/health", healthHandler)

	// cache refresh endpoint
	router.GET("/refresh", refreshHandler)

	// commands endpoint
	router.POST(slackRouterPrefix("commands"), commandsHandler)

	// events endpoint
	router.POST(slackRouterPrefix("events"), eventsHandler)

	// interactions endpoint
	router.POST(slackRouterPrefix("interactions"), interactionsHandler)

	var tomb tomb.Tomb
	// start http server
	tomb.Go(func() error {
		srv := &http.Server{Addr: ":" + port, Handler: router}

		go func() {
			logrus.Info("Starting HTTP server...")
			err := srv.ListenAndServe()
			if err != nil {
				tomb.Kill(err)
			}
		}()

		for {
			select {
			case <-tomb.Dying():
				logrus.Info("Stopping HTTP server...")
				return srv.Shutdown(context.Background())
			}
		}
	})

	// watch for errors and terminate safely
	tomb.Wait()

	return tomb.Err()
}

// slackRouterPrefix returns appropriate an optional router prefix
func slackRouterPrefix(endpoint string) string {

	// grab prefix from the env
	prefix := strings.ToLower(os.Getenv("slack_router_prefix"))

	if prefix != "" {
		return prefix + endpoint
	}

	// set to / as a default
	return "/" + endpoint
}
