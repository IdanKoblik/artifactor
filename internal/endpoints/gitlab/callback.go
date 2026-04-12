package gitlab

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"packster/internal/repository"
	"packster/internal/requests"
	"packster/internal/utils"
	"packster/pkg/types"
	gitlabtypes "packster/pkg/types/gitlab"

	"github.com/gin-gonic/gin"
)

func (h *GitlabHandler) HandleCallback(c *gin.Context) {
	state := c.Query("state")
	hostId, err := strconv.Atoi(state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	host, err := utils.GetHostById(h.DB, hostId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("host %d not found", hostId)})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	redirectURI := fmt.Sprintf("%s://%s/api/auth/gitlab/callback", scheme, c.Request.Host)
	payload := map[string]string{
		"client_id":     host.ApplicationId,
		"client_secret": host.Secret,
		"code":          c.Query("code"),
		"grant_type":    "authorization_code",
		"redirect_uri":  redirectURI,
	}

	type oauthResult struct {
		token *gitlabtypes.OauthToken
		err   error
	}

	ch := make(chan oauthResult, 1)
	client := &http.Client{}

	go func() {
		res, err := requests.GitlabOauthToken(client, payload, host.Url)
		ch <- oauthResult{token: res, err: err}
	}()

	result := <-ch
	if result.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.err.Error(),
		})
		return
	}

	user, err := requests.FetchGitlabUser(client, result.token.Token, host.Url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	accountReq := types.AuthRequest{
		Username: user.Username,
		SsoId:    strconv.Itoa(user.ID),
		Host:     user.Host,
	}

	account, err := h.Repo.CreateAccount(accountReq)
	if errors.Is(err, repository.ErrAccountExists) {
		_, err = h.Repo.GetDB().Exec(`UPDATE account SET last_login=$1 WHERE id=(SELECT account FROM auth WHERE username=$2 AND sso_id=$3 AND host=$4)`,
			time.Now(),
			account.AuthData.Username,
			account.AuthData.SsoId,
			account.AuthData.Host,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, user)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
