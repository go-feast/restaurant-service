package handler

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/topics"
)

func (h *Handler) ReceiveOrderEvent(topic topics.Topic) message.NoPublishHandlerFunc {
	const format = "Received %s event"

	var str = fmt.Sprintf(format, topic)

	return func(_ *message.Message) error {
		h.logger.Info().Msg(str)
		return nil
	}
}
