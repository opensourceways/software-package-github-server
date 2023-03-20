package domain

type PullRequest interface {
	CreateRepo(repo string) (string, error)
}
