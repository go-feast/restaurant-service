package handler

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"service/domain/restaurant"
)

func (h *Handler) CreateOrder(msg *message.Message) ([]*message.Message, error) {
	var (
		ctx       = msg.Context()
		orderPaid restaurant.OrderPaid
	)

	err := h.unmarshaler.Unmarshal(msg.Payload, &orderPaid)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling order paid")
	}

	rest, err := h.restaurantRepository.Get(ctx, orderPaid.RestaurantID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting restaurant from repository")
	}

	meals := mapMeals(rest.Meals, orderPaid.Meals)

	order := restaurant.NewOrder(
		orderPaid.OrderID,
		orderPaid.RestaurantID,
		meals,
	)

	err = h.orderRepository.Create(ctx, order)
	if err != nil {
		return nil, errors.Wrap(err, "error creating order")
	}

	return []*message.Message{}, nil
}

func mapMeals(searchFrom restaurant.Meals, ids uuid.UUIDs) restaurant.Meals {
	mapMeals := make(map[uuid.UUID]*restaurant.Meal)

	for _, meal := range searchFrom {
		mapMeals[meal.ID] = meal
	}

	newMeals := make(restaurant.Meals, len(ids), len(ids)) //nolint:gosimple

	for i, id := range ids {
		newMeals[i] = mapMeals[id]
	}

	return newMeals
}
