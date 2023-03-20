package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v36/github"

	"github.com/opensourceways/software-package-github-server/mq"
)

const platformGithub = "github"

type iClient interface {
	GetRepo(org, repo string) (*github.Repository, error)
	CreateRepo(org string, r *github.Repository) error
}

func newMessageService(cli iClient, cfg *Config) *message {
	return &message{
		cli: cli,
		cfg: cfg,
	}
}

type message struct {
	cli iClient
	cfg *Config
}

func (m *message) subscribe() error {
	h := map[string]mq.Handler{
		m.cfg.Topics.ApprovedPkg: m.handleCreateRepo,
		m.cfg.Topics.MergedPR:    m.handleCreateRepo,
	}

	return mq.Subscriber().Subscribe(m.cfg.MQ.GroupName, h)
}

type CreateRepoEvent struct {
	PkgId   string `json:"pkg_id"`
	PkgName string `json:"pkg_name"`
}

func (m *message) handleCreateRepo(data []byte) error {
	msg := new(CreateRepoEvent)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	if err := m.createRepo(msg.PkgName); err != nil {
		return err
	}

	return m.notifyCreatedRepo(msg.PkgId, msg.PkgName)
}

func (m *message) createRepo(repo string) error {
	v, err := m.cli.GetRepo(m.cfg.Org, repo)
	if err == nil && *v.Name == repo {
		return nil
	}

	r := &github.Repository{Name: &repo}

	return m.cli.CreateRepo(m.cfg.Org, r)
}

type CreatedRepoMsg struct {
	PkgId    string `json:"pkg_id"`
	Platform string `json:"platform"`
	RepoLink string `json:"repo_link"`
}

func (m *message) notifyCreatedRepo(pkgId, pkgName string) error {
	link := fmt.Sprintf("https://github.com/%s/%s", m.cfg.Org, pkgName)
	e := CreatedRepoMsg{
		PkgId:    pkgId,
		Platform: platformGithub,
		RepoLink: link,
	}

	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return mq.Subscriber().Publish(m.cfg.Topics.CreatedRepo, body)
}
