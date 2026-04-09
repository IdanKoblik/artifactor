package gitlab

import (
	"packster/internal/repository"
	"packster/pkg/config"
)

type GitlabHandler struct {
	Repo repository.IAccountRepo
	Cfg *config.Config
}

func NewGitlabHandler(repo *repository.AccountRepo, cfg *config.Config) *GitlabHandler {
	return &GitlabHandler{
		Repo: repo,
		Cfg: cfg,
	}
}
