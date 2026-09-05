package middleware

import (
	"github.com/QuantumNous/new-api/service"

	"github.com/gin-gonic/gin"
)

// DatabaseActivityMaintenance runs due maintenance after a real API or relay
// request. Static assets and the in-memory health endpoint do not wake the
// database merely to perform maintenance.
func DatabaseActivityMaintenance() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		c.Next()
		if shouldNotifyDatabaseActivity(path, c.GetString(RouteTagKey)) {
			service.NotifyDatabaseActivity()
		}
	}
}

func shouldNotifyDatabaseActivity(path string, routeTag string) bool {
	if path == "/api/status" {
		return false
	}
	if routeTag == "api" || routeTag == "old_api" || routeTag == "relay" {
		return true
	}
	return false
}
