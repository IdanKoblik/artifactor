package endpoints

import (
	"net/http"
    "database/sql"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	PgsqlStatus string `json:"pgsql"`
}

// HandleHealth godoc
// @Summary      Health check
// @Description  Returns the health status of required services.
// @Tags         system
// @Produce      json
// @Success      200  {object}  types.HealthResponse  "All services healthy"
// @Failure      500  {object}  types.HealthResponse  "One or more services unhealthy"
// @Router       /api/health [get]
func HandleHealth(c *gin.Context, sqlConn *sql.DB) {
	response := HealthResponse{
		PgsqlStatus: "Pgsql is fine",
	}

	status := http.StatusOK
	if sqlConn == nil {
		status = http.StatusInternalServerError
		response.PgsqlStatus = "Pgsql instance was not found"
	} else {
		err := sqlConn.Ping()
		if err != nil {
			status = http.StatusInternalServerError
			response.PgsqlStatus = err.Error()
		}
	}

	c.JSON(status, response)
}
