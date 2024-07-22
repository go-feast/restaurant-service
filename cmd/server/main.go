package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service/api/http/restaurant/handler"
	"service/closer"
	"service/config"
	"service/domain/restaurant"
	"service/event"
	mw "service/http/middleware"
	"service/infrastructure/repository/order"
	rest "service/infrastructure/repository/restaurant"
	"service/infrastructure/saver"
	"service/logging"
	"service/metrics"
	"service/pubsub"
	serv "service/server"
	"service/tracing"
)

const (
	version     = "v1.0"
	serviceName = "restaurant"
	driverName  = "pgx/v5"
)

func main() {
	c := &config.ServiceConfig{}
	// config
	err := config.ParseConfig(c)
	if err != nil {
		log.Fatal(err)
	}
	// logger
	logger := logging.New(
		logging.WithTimestamp(),
		logging.WithServiceName(serviceName),
		logging.WithPID(),
	)

	logger.Info().Any("config", c).Send()

	forClose := closer.NewCloser(logger)
	defer forClose.Close()

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(version),
		),
		resource.WithProcess(),
		resource.WithOS(),
	)
	if err != nil {
		logger.Err(err).Msg("filed to create resource")
		return
	}

	if err = tracing.RegisterTracerProvider(ctx, res); err != nil {
		logger.Err(err).Msg("failed to register tracer provider")
		return
	}

	metrics.RegisterServiceName(serviceName)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName:           "pgx",
		DSN:                  c.DB.DSN(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
		return
	}

	db = db.WithContext(ctx)

	// domain gorm initialization
	err = restaurant.InitRestaurantModels(db)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize restaurant models")
		return
	}

	// main server
	mainServiceServer, mainRouter := serv.NewServer(c.Server)

	// metric server
	metricServer, metricRouter := serv.NewServer(c.MetricServer)

	// register routes
	//		main
	fc := RegisterMainServiceRoutes(mainRouter, db, *c.Kafka)

	forClose.AppendClosers(fc...)
	//		metric
	RegisterMetricRoute(metricRouter)

	_, errCh := serv.Run(ctx, mainServiceServer, metricServer)

	for err = range errCh {
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Err(err).Send()
		}
	}
}

func Middlewares(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(logging.NewLogEntry()))
	r.Use(middleware.Recoverer)
}

func RegisterMainServiceRoutes(
	r chi.Router,
	db *gorm.DB,
	c config.KafkaConfig,
) []closer.C { //nolint:unparam
	// middlewares
	Middlewares(r)
	r.Get("/healthz", mw.Healthz)

	kafkaPub, err := pubsub.NewKafkaPublisher(c.KafkaURL, logging.NewWatermillAdapter())
	if err != nil {
		panic(err)
	}

	restaurantRepository := rest.NewGormRepository(db)
	orderRepository := order.NewGormRepository(db)
	s := saver.NewSaver(restaurantRepository)
	h := handler.NewHandler(s, orderRepository, kafkaPub, event.JSONMarshaler{}, logging.New())

	r.Route("/api/v1/", func(r chi.Router) {
		r.Route("restaurants/", func(r chi.Router) {
			r.Post("/", h.CreateRestaurant) // create rest
		})
		r.Route("/orders", func(r chi.Router) {
			r.Post("/{id}", h.CookingOrder)
		})
	})

	return []closer.C{}
}

func RegisterMetricRoute(r chi.Router) {
	h := promhttp.Handler()
	r.Get("/metrics", h.ServeHTTP)
}
