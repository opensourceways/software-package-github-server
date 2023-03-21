package domain

type Repo interface {
	CreateRepo(repo string) (string, error)
}
