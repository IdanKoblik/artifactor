package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func InternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func IsAdmin(c *gin.Context) bool {
	return c.GetBool("admin")
}

func GetToken(c *gin.Context) string {
	return c.GetString("token")
}
