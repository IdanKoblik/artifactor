package gitlab

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"packster/internal/utils"
	"packster/pkg/types"

	"github.com/gin-gonic/gin"
)

func (h *GitlabHandler) HandleRedirect(c *gin.Context) {
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid host id"})
		return
	}

	host, err := utils.GetHostById(h.DB, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("host %d not found", id)})
		return
	}

	if host.Type != types.Gitlab {
		c.JSON(http.StatusBadRequest, gin.H{"error": "host is not a gitlab instance"})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	redirectURL := buildRedirectUrl(host, scheme, c.Request.Host)
	c.Redirect(http.StatusFound, redirectURL)
}

func buildRedirectUrl(host *types.Host, scheme, reqHost string) string {
	baseURL := fmt.Sprintf("%s/oauth/authorize", host.Url)
	redirectURI := fmt.Sprintf("%s://%s/api/auth/gitlab/callback", scheme, reqHost)

	params := url.Values{}
	params.Add("client_id", host.ApplicationId)
	params.Add("redirect_uri", redirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "read_user")
	params.Add("state", strconv.Itoa(host.Id))

	return baseURL + "?" + params.Encode()
}
