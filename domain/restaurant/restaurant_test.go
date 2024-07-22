package restaurant_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"service/domain/restaurant"
	"testing"
)

func TestMeals_Append(t *testing.T) {
	meal := restaurant.Meal{ID: uuid.New()}

	meals := restaurant.NewMeals(uint(0))

	meals.Append(meal)

	assert.Contains(t, meals, &meal)
}
