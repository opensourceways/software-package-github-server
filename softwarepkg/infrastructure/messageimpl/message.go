package messageimpl

import (
	"github.com/opensourceways/software-package-github-server/message-server"
	"github.com/opensourceways/software-package-github-server/mq"
	"github.com/opensourceways/software-package-github-server/softwarepkg/domain"
)

func NewMessageImpl(t messageserver.TopicsToNotify) *MessageImpl {
	return &MessageImpl{
		topics: t,
	}
}

type MessageImpl struct {
	topics messageserver.TopicsToNotify
}

func (m *MessageImpl) NotifyRepoCreatedResult(msg domain.EventMessage) error {
	return send(m.topics.CreatedRepo, msg)
}

func send(topic string, v domain.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	return mq.Subscriber().Publish(topic, body)
}
