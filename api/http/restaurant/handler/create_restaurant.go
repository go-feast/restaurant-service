package handler

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"net/http"
	"service/domain/restaurant"
	"service/domain/shared/location"
	"service/http/httpstatus"
)

type CreateRestaurantRequest struct { //nolint:govet
	Name               string
	ContactInformation map[string]string
	location           location.JSONLocation
}

type CreateRestaurantResponse struct { //nolint:govet
	Name               string
	ContactInformation map[string]string
	location           location.JSONLocation
}

func (h *Handler) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var (
		createRestaurantRequest CreateRestaurantRequest
		ctx                     = r.Context()
	)

	err := json.NewDecoder(r.Body).Decode(&createRestaurantRequest)
	if err != nil {
		httpstatus.BadRequest(ctx, w, errors.Wrap(err, "invalid body"))
		return
	}

	rest, err := restaurant.NewRestaurant(
		createRestaurantRequest.Name,
		createRestaurantRequest.ContactInformation,
		createRestaurantRequest.location.ToLocation(),
		nil,
	)
	if err != nil {
		httpstatus.BadRequest(ctx, w, errors.Wrap(err, "failed to create restaurant"))
		return
	}

	err = h.saver.Save(ctx, rest)
	if err != nil {
		httpstatus.InternalServerError(ctx, w, errors.Wrap(err, "failed to save restaurant"))
		return
	}

	response := CreateRestaurantResponse{
		Name:               rest.Name,
		ContactInformation: rest.ContactInformation,
		location:           rest.GeoPosition.ToJSON(),
	}

	render.JSON(w, r, response)
}
