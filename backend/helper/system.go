package helper

import (
	"hbd/structs"

	"github.com/gin-gonic/gin"
)

// HealthCheck checks if the service is ready and returns a status response.
// @Summary Check service readiness
// @Description This endpoint checks the readiness of the service and returns a status.
// @Produce  json
// @Success 200 {object} structs.Ready
// @Router /health [get]
// @Tags health
func HealthCheck(c *gin.Context) {
	c.JSON(200, structs.Ready{Status: "ok"})
}
