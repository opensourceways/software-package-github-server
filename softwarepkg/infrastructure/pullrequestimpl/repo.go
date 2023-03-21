package pullrequestimpl

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v36/github"
)

type iClient interface {
	GetRepo(org, repo string) (*github.Repository, error)
	CreateRepo(org string, r *github.Repository) error
}

func NewRepoImpl(cfg Config, cli iClient) *RepoImpl {
	return &RepoImpl{
		cfg: cfg,
		cli: cli,
	}
}

type RepoImpl struct {
	cfg Config
	cli iClient
}

func (impl *RepoImpl) CreateRepo(repo string) (string, error) {
	r := &github.Repository{Name: &repo}
	err := impl.cli.CreateRepo(impl.cfg.Org, r)
	if err != nil && !strings.Contains(err.Error(), "name already exists") {
		return "", err
	}

	url := fmt.Sprintf("https://github.com/%s/%s", impl.cfg.Org, repo)

	return url, nil
}
