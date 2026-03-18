package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleFetch godoc
// @Summary      Fetch an API token
// @Description  Returns metadata for the specified token. Requires admin privileges.
// @Tags         auth
// @Produce      json
// @Param        token  path  string  true  "Token to look up"
// @Success      200  {object}  types.ApiToken  "Token metadata"
// @Failure      400  {object}  map[string]string  "Token not found or missing"
// @Failure      401  {string}  string  "Admin privileges required"
// @Security     ApiKeyAuth
// @Router       /fetch/{token} [get]
func (h *AuthHandler) HandleFetch(c *gin.Context) {
	admin, exists := c.Get("admin")
	if !exists || !admin.(bool) {
		c.String(http.StatusUnauthorized, "Only admin allowed to prune an api token")
		return
	}

	rawToken := c.Param("token")
	if rawToken == "" {
		c.String(http.StatusBadRequest, "Missing api token")
		return
	}

	token, err := h.Repo.FetchToken(rawToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, token)
}
