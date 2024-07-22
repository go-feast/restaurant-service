package restaurant

import "github.com/google/uuid"

type Order struct { //nolint:govet
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	RestaurantID uuid.UUID `gorm:"type:uuid"`
	Meals        Meals     `gorm:"references:RestaurantID"`
	Cooking      bool
	Finished     bool
}

func NewOrder(
	id uuid.UUID,
	restaurantID uuid.UUID,
	meals Meals,
) *Order {
	return &Order{
		ID:           id,
		RestaurantID: restaurantID,
		Meals:        meals,
		Cooking:      false,
		Finished:     false,
	}
}
