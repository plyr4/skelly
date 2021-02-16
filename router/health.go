package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthHandler represents the API handler to
// report the health status for the application.
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
