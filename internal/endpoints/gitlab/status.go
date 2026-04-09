package gitlab

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *GitlabHandler) HandleStatus(c *gin.Context) {
	if h.Cfg.Gitlab == nil {
		c.String(http.StatusBadRequest, "Gitlab is not enabled")
		return
	}

	c.Status(200)
}
