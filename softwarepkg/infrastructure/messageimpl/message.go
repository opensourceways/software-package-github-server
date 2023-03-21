package messageimpl

import (
	"github.com/opensourceways/software-package-github-server/mq"
	"github.com/opensourceways/software-package-github-server/softwarepkg/domain"
)

func NewMessageImpl(c Config) *MessageImpl {
	return &MessageImpl{
		cfg: c,
	}
}

type MessageImpl struct {
	cfg Config
}

func (m *MessageImpl) NotifyRepoCreatedResult(msg domain.EventMessage) error {
	return send(m.cfg.TopicsToNotify.CreatedRepo, msg)
}

func send(topic string, v domain.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return mq.Subscriber().Publish(topic, body)
}
