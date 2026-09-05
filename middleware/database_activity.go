package middleware

import (
	"strings"

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
	if path == "/api" || strings.HasPrefix(path, "/api/") ||
		path == "/v1" || strings.HasPrefix(path, "/v1/") ||
		path == "/v1beta" || strings.HasPrefix(path, "/v1beta/") ||
		path == "/pg" || strings.HasPrefix(path, "/pg/") ||
		path == "/mj" || strings.HasPrefix(path, "/mj/") {
		return true
	}
	return strings.Contains(path, "/mj/")
}
