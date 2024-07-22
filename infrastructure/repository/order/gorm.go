package order

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

func (g *GormRepository) Create(ctx context.Context, o *restaurant.Order) error {
	return g.db.WithContext(ctx).Create(o).Error
}

func (g *GormRepository) Get(ctx context.Context, id uuid.UUID) (*restaurant.Order, error) {
	var o restaurant.Order

	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Find(&o, "id = ?", id).Error
		if err != nil {
			return errors.Wrap(err, "failed to find order")
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "transaction failed")
	}

	return &o, nil
}

func (g *GormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return g.db.WithContext(ctx).Delete(&restaurant.Order{}, "id = ?", id).Error
}
