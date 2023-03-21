package app

import (
	"github.com/opensourceways/software-package-github-server/softwarepkg/domain"
)

type MessageService interface {
	CreateRepo(CmdToCreateRepo) error
}

func NewMessageService(
	p domain.Repo,
	s domain.SoftwarePkgProducer,
) *messageService {
	return &messageService{
		pr:       p,
		producer: s,
	}
}

type messageService struct {
	pr       domain.Repo
	producer domain.SoftwarePkgProducer
}

func (m *messageService) CreateRepo(cmd CmdToCreateRepo) error {
	url, err := m.pr.CreateRepo(cmd.PkgName)
	if err != nil {
		return err
	}

	e := domain.NewRepoCreatedEvent(cmd.PkgId, url)
	return m.producer.NotifyRepoCreatedResult(&e)
}
