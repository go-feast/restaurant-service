package restaurant

import (
	"github.com/google/uuid"
	"service/domain/shared/location"
)

type Currency string

const (
	USD Currency = "USD" //nolint:revive
)

type Restaurant struct { //nolint:govet
	ID                 uuid.UUID
	Name               string
	ContactInformation map[string]string
	GeoPosition        location.Location
	Meals              Meals
}

func NewRestaurant(
	name string,
	contactInformation map[string]string,
	l location.Location,
	meals Meals,
) (*Restaurant, error) {
	return &Restaurant{
		ID:                 uuid.New(),
		Name:               name,
		ContactInformation: contactInformation,
		GeoPosition:        l,
		Meals:              meals,
	}, nil
}

type Meal struct { //nolint:govet
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	Currency    Currency
}

func NewMeal(
	name string,
	description string,
	price float64,
) *Meal {
	return &Meal{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
	}
}

type Meals []*Meal

func NewMeals(n uint) Meals {
	return make(Meals, 0, n)
}

func (m *Meals) Append(meal Meal) {
	*m = append(*m, &meal)
}
