package auth

import (
	"net/http"
	"packster/internal/endpoints"

	"github.com/gin-gonic/gin"
)

// HandleFetch godoc
// @Summary      Fetch an API token
// @Description  Returns metadata for the specified token.
// @Tags         auth
// @Produce      json
// @Param        token  path  string  true  "Token to look up"
// @Success      200  {object}  types.ApiToken  "Token metadata"
// @Failure      400  {object}  map[string]string  "Token not found or missing"
// @Security     ApiKeyAuth
// @Router       /fetch/{token} [get]
func (h *AuthHandler) HandleFetch(c *gin.Context) {
	rawToken := c.Param("token")
	if rawToken == "" {
		c.String(http.StatusBadRequest, "Missing api token")
		return
	}

	token, err := h.Repo.FetchToken(rawToken)
	if err != nil {
		endpoints.BadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, token)
}
