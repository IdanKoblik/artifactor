package endpoints

import (
	"database/sql"
	"net/http"

	"packster/internal/utils"

	"github.com/gin-gonic/gin"
)

type HostEntry struct {
	Id   int    `json:"id"`
	Url  string `json:"url"`
	Type string `json:"type"`
}

func HandleHosts(c *gin.Context, sqlConn *sql.DB) {
	hosts, err := utils.GetHosts(sqlConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]HostEntry, 0, len(hosts))
	for _, h := range hosts {
		result = append(result, HostEntry{
			Id:   h.Id,
			Url:  h.Url,
			Type: h.Type.String(),
		})
	}

	c.JSON(http.StatusOK, result)
}
