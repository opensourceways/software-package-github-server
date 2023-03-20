package pullrequestimpl

import (
	"fmt"

	"github.com/google/go-github/v36/github"

	"github.com/opensourceways/software-package-github-server/config"
)

type iClient interface {
	GetRepo(org, repo string) (*github.Repository, error)
	CreateRepo(org string, r *github.Repository) error
}

func NewPullRequestImpl(cfg *config.Config, cli iClient) *PullRequestImpl {
	return &PullRequestImpl{
		cfg: cfg,
		cli: cli,
	}
}

type PullRequestImpl struct {
	cfg *config.Config
	cli iClient
}

func (impl *PullRequestImpl) CreateRepo(repo string) (url string, err error) {
	org := impl.cfg.Org
	url = fmt.Sprintf("https://github.com/%s/%s", org, repo)

	_, err = impl.cli.GetRepo(org, repo)
	if err == nil {

		return
	}

	r := &github.Repository{Name: &repo}
	err = impl.cli.CreateRepo(org, r)

	return
}
