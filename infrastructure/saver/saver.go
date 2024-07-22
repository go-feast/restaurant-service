package saver

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"service/domain/restaurant"
	"service/event"
)

type Saver struct {
	_          message.Publisher
	_          event.Marshaler
	repository restaurant.Repository[*restaurant.Restaurant]
}

func NewSaver(repository restaurant.Repository[*restaurant.Restaurant]) *Saver {
	return &Saver{repository: repository}
}

func (s Saver) Save(ctx context.Context, r *restaurant.Restaurant) error {
	err := s.repository.Create(ctx, r)
	if err != nil {
		return errors.Wrap(err, "saving: failed to create order")
	}

	return nil
}
