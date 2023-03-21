package domain

import "encoding/json"

const platformGithub = "github"

type RepoCreatedEvent struct {
	PkgId    string `json:"pkg_id"`
	Platform string `json:"platform"`
	RepoLink string `json:"repo_link"`
}

func (e *RepoCreatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

func NewRepoCreatedEvent(pkgId, url string) RepoCreatedEvent {
	return RepoCreatedEvent{
		PkgId:    pkgId,
		Platform: platformGithub,
		RepoLink: url,
	}
}
