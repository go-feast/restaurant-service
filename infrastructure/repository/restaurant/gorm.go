package restaurant

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"service/domain/restaurant"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (g *GormRepository) Create(ctx context.Context, r *restaurant.Restaurant) error {
	return g.db.WithContext(ctx).Create(r.ToDatabaseDTO()).Error
}

func (g *GormRepository) Get(ctx context.Context, id uuid.UUID) (*restaurant.Restaurant, error) {
	var r restaurant.DatabaseRestaurantDTO

	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Find(&r, id).Error
		if err != nil {
			return errors.Wrap(err, "failed to find restaurant")
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "transaction failed")
	}

	return r.ToRestaurant(), nil
}

func (g *GormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		r, err := withTx(tx).Get(ctx, id)
		if err != nil {
			return errors.Wrap(err, "failed to find restaurant")
		}

		err = tx.Delete(r).Error
		if err != nil {
			return errors.Wrap(err, "failed to delete restaurant")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}

func withTx(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}
