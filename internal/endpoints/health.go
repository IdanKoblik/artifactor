package endpoints

import (
	"net/http"
    "database/sql"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	PgsqlStatus string `json:"pgsql"`
}

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
