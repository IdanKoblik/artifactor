package endpoints

import (
	"net/http"
	"artifactor/internal/sql"
	"artifactor/internal/redis"
	responses "artifactor/pkg/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) HandleHealth(c *gin.Context) {
	response := responses.HealthResponse{
		SqlStatus: "Sql is fine",
		RedisStatus: "Redis is fine",
	}

	status := http.StatusOK
	err := sql.CheckHealth()
	if err != nil {
		response.SqlStatus = err.Error()
		status = http.StatusInternalServerError
	}

	err = redis.CheckHealth()
	if err != nil {
		response.RedisStatus = err.Error()
		status = http.StatusInternalServerError
	}

	c.JSON(status, response)
}
