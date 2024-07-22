package handler

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog"
	"service/domain/restaurant"
	"service/domain/shared/saver"
	"service/event"
)

type Handler struct {
	saver           saver.Saver[*restaurant.Restaurant]
	orderRepository restaurant.Repository[*restaurant.Order]
	publisher       message.Publisher
	marshaller      event.Marshaler
	logger          *zerolog.Logger
}

func NewHandler(
	s saver.Saver[*restaurant.Restaurant],
	orderRepository restaurant.Repository[*restaurant.Order],
	publisher message.Publisher,
	marshaller event.Marshaler,
	logger *zerolog.Logger,
) *Handler {
	return &Handler{
		saver:           s,
		orderRepository: orderRepository,
		publisher:       publisher,
		marshaller:      marshaller,
		logger:          logger,
	}
}
