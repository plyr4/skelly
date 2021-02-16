package router

import (
	"net/http"

	"github.com/davidvader/skelly/emojis"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// refreshHandler represents the API handler to
// automatically refresh any cache loaded into memory.
func refreshHandler(c *gin.Context) {

	logrus.Info("manually refreshing application cache")

	// refresh emojis cache
	err := emojis.Load()
	if err != nil {
		err = errors.Wrap(err, "could not refresh emojis cache")
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	// refresh successful
	logrus.Info("manual refresh successful")

	c.JSON(http.StatusOK, "manual refresh successful")
}
