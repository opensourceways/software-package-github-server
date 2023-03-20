package messageserver

import (
	"encoding/json"

	"github.com/opensourceways/software-package-github-server/mq"
	"github.com/opensourceways/software-package-github-server/softwarepkg/app"
)

func Init(s app.MessageService) *MessageServer {
	return &MessageServer{
		service: s,
	}
}

type MessageServer struct {
	service app.MessageService
}

func (m *MessageServer) Subscribe(cfg *Config) error {
	h := map[string]mq.Handler{
		cfg.Topics.ApprovedPkg: m.handleCreateRepo,
		cfg.Topics.MergedPR:    m.handleCreateRepo,
	}

	return mq.Subscriber().Subscribe(cfg.Group, h)
}

func (m *MessageServer) handleCreateRepo(data []byte) error {
	msg := new(CreateRepoEvent)

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	return m.service.HandleCreateRepo(*msg)
}
