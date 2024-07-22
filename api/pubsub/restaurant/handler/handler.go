package handler

import (
	"github.com/rs/zerolog"
	"service/domain/restaurant"
	"service/event"
)

type Handler struct {
	unmarshaler          event.Unmarshaler
	restaurantRepository restaurant.Repository[*restaurant.Restaurant]
	orderRepository      restaurant.Repository[*restaurant.Order]
	logger               *zerolog.Logger
}

func NewHandler(
	unmarshaler event.Unmarshaler,
	restaurantRepository restaurant.Repository[*restaurant.Restaurant],
	orderRepository restaurant.Repository[*restaurant.Order],
	logger *zerolog.Logger,
) *Handler {
	return &Handler{
		unmarshaler:          unmarshaler,
		restaurantRepository: restaurantRepository,
		orderRepository:      orderRepository,
		logger:               logger,
	}
}
