package gitlab

import (
	"database/sql"

	"packster/internal/repository"
)

type GitlabHandler struct {
	Repo repository.IAccountRepo
	DB   *sql.DB
}

func NewGitlabHandler(repo repository.IAccountRepo, db *sql.DB) *GitlabHandler {
	return &GitlabHandler{
		Repo: repo,
		DB:   db,
	}
}
