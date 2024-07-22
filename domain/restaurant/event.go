package restaurant

import (
	"github.com/google/uuid"
	"service/event"
)

type OrderPaid struct { //nolint:govet
	event.Event  `json:"-"`
	OrderID      uuid.UUID  `json:"order_id"`
	RestaurantID uuid.UUID  `json:"restaurant_id"`
	Meals        uuid.UUIDs `json:"meals"`
}

type OrderCooking struct {
	event.Event `json:"-"`
	OrderID     uuid.UUID `json:"order_id"`
}

type OrderFinished struct {
	event.Event `json:"-"`
	OrderID     uuid.UUID `json:"order_id"`
}
