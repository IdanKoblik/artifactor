package requests

import (
	"fmt"
	"bytes"
	"io"
	"encoding/json"
	"net/http"

	"packster/pkg/types/gitlab"
)

func GitlabOauthToken(client *http.Client, payload map[string]string, gitlabHost string) (*gitlab.OauthToken, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", gitlabHost+"/oauth/token", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res gitlab.OauthToken
	_ = json.Unmarshal(respBody, &res)
	return &res, nil
}

func FetchGitlabUser(client *http.Client, token, gitlabHost string) (*gitlab.GitlabUser, error) {
	req, err := http.NewRequest("GET", gitlabHost+"/api/v4/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res gitlab.GitlabUser
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if res.ID == 0 {
		return nil, fmt.Errorf("failed to get Gitlab user")
	}

	res.Host = gitlabHost
	res.Token = token
	return &res, nil
}
