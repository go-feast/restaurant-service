package handler

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"
	"github.com/go-feast/topics"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"service/domain/restaurant"
	"service/event"
	"service/http/httpstatus"
)

func (h *Handler) CookingOrder(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpstatus.BadRequest(ctx, w, errors.Wrap(err, "invalid id"))
		return
	}

	//ERROR: unhandled behavior when sent second reqest on same order
	order, err := h.orderRepository.Get(ctx, id)
	if err != nil {
		return
	}

	var (
		e     event.Event
		topic topics.Topic
	)

	switch {
	case !order.Cooking:
		e = restaurant.OrderCooking{
			OrderID: order.ID,
		}
		topic = topics.Cooking
	case !order.Finished:
		e = restaurant.OrderFinished{
			OrderID: order.ID,
		}
		topic = topics.CookingFinished
	}

	marshal, err := h.marshaller.Marshal(e)
	if err != nil {
		httpstatus.InternalServerError(ctx, w, errors.Wrap(err, "marshaller failed"))
		return
	}

	m := message.NewMessage(uuid.NewString(), marshal)

	err = h.publisher.Publish(topic.String(), m)
	if err != nil {
		httpstatus.InternalServerError(ctx, w, errors.Wrap(err, "publish failed"))
		return
	}

	h.logger.Info().Str("order_id", id.String()).Msgf("order state changed to %s", topic.String())

	w.WriteHeader(http.StatusOK)
}
