package restaurant

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"service/domain/shared/location"
)

func InitRestaurantModels(db *gorm.DB) error {
	return db.AutoMigrate(&DatabaseRestaurantDTO{}, &DatabaseMealDTO{}, &Order{})
}

// DatabaseRestaurantDTO represents the GORM DTO for the Restaurant struct.
type DatabaseRestaurantDTO struct { //nolint:govet
	ID                 uuid.UUID         `gorm:"type:uuid;primaryKey"`
	Name               string            `gorm:"type:text"`
	ContactInformation map[string]string `gorm:"type:jsonb"`
	LocationLatitude   float64           `gorm:"type:numeric"`
	LocationLongitude  float64           `gorm:"type:numeric"`
	Meals              []DatabaseMealDTO `gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`
	Order              []Order           `gorm:"foreignKey:RestaurantID;references:ID;constraint:OnDelete:CASCADE"`
}

// DatabaseMealDTO represents the GORM DTO for the Meal struct.
type DatabaseMealDTO struct { //nolint:govet
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	RestaurantID uuid.UUID `gorm:"type:uuid"`
	Name         string    `gorm:"type:text"`
	Description  string    `gorm:"type:text"`
	Price        float64   `gorm:"type:numeric"`
	Currency     Currency  `gorm:"type:char(3)"`
}

// ToRestaurant converts a DatabaseRestaurantDTO to a Restaurant.
func (d *DatabaseRestaurantDTO) ToRestaurant() *Restaurant {
	loc := location.NewLocation(d.LocationLatitude, d.LocationLongitude)

	meals := make(Meals, len(d.Meals))
	for i, meal := range d.Meals {
		meals[i] = meal.ToMeal()
	}

	return &Restaurant{
		ID:                 d.ID,
		Name:               d.Name,
		ContactInformation: d.ContactInformation,
		GeoPosition:        loc,
		Meals:              meals,
	}
}

// ToMeal converts a DatabaseMealDTO to a Meal.
func (d *DatabaseMealDTO) ToMeal() *Meal {
	return &Meal{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		Currency:    d.Currency,
	}
}

// ToDatabaseDTO converts a Restaurant to a DatabaseRestaurantDTO.
func (r *Restaurant) ToDatabaseDTO() *DatabaseRestaurantDTO {
	meals := make([]DatabaseMealDTO, len(r.Meals))
	for i, meal := range r.Meals {
		meals[i] = meal.ToDatabaseDTO(r.ID)
	}

	return &DatabaseRestaurantDTO{
		ID:                 r.ID,
		Name:               r.Name,
		ContactInformation: r.ContactInformation,
		LocationLatitude:   r.GeoPosition.Latitude(),
		LocationLongitude:  r.GeoPosition.Longitude(),
		Meals:              meals,
	}
}

// ToDatabaseDTO converts a Meal to a DatabaseMealDTO.
func (m *Meal) ToDatabaseDTO(restaurantID uuid.UUID) DatabaseMealDTO {
	return DatabaseMealDTO{
		ID:           m.ID,
		RestaurantID: restaurantID,
		Name:         m.Name,
		Description:  m.Description,
		Price:        m.Price,
		Currency:     m.Currency,
	}
}
