//go:generate oapi-codegen --config=../../internal/worker/codegen/spec.yaml ../../doc/worker/worker.yaml
//go:generate oapi-codegen --config=../../internal/worker/codegen/types.yaml ../../doc/worker/worker.yaml
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"

	"github.com/AddMile/backend/internal/app/user"
	"github.com/AddMile/backend/internal/config"
	userJob "github.com/AddMile/backend/internal/job/user"
	gen "github.com/AddMile/backend/internal/worker/codegen"
	worker "github.com/AddMile/backend/internal/worker/handler"

	codegenkit "github.com/AddMile/backend/internal/kit/codegen"
	ciokit "github.com/AddMile/backend/internal/kit/customerio"
	httpserverkit "github.com/AddMile/backend/internal/kit/httpserver"
	loggerkit "github.com/AddMile/backend/internal/kit/logger"
	pgkit "github.com/AddMile/backend/internal/kit/pg"
	pubsubkit "github.com/AddMile/backend/internal/kit/pubsub"
	validatorkit "github.com/AddMile/backend/internal/kit/validator"
)

type Worker struct {
	*worker.UserHTTPHandler
}

func NewWorker(
	userHandler *worker.UserHTTPHandler,
) Worker {
	return Worker{
		userHandler,
	}
}

func main() {
	ctx := context.Background()
	cfg := config.Load()

	l := loggerkit.NewLogger(os.Stdout, cfg.Debug)

	validator := validatorkit.New()

	db, err := pgkit.New(cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}

	publisher, err := pubsubkit.NewPublisher(ctx, cfg.GoogleProjectID)
	if err != nil {
		log.Fatal(err)
	}

	cioConfig := ciokit.Config{
		APIKey:    cfg.CustomerIOAPIKey,
		Endpoint:  cfg.CustomerIOEndpoint,
		BatchSize: cfg.CustomerIOBatchSize,
		Interval:  time.Millisecond * time.Duration(cfg.CustomerIOFlushInterval),
		Verbose:   cfg.CustomerIOVerbose,
	}
	cioClient, err := ciokit.NewClient(cioConfig)
	if err != nil {
		log.Fatalf("creating customer.io client: %v", err)
	}

	userStorage := user.NewStorage(db)
	userQueue := user.NewQueue(cfg, publisher)
	userService, err := user.NewService(l, validator, userStorage, userQueue)
	if err != nil {
		log.Fatalf("initializing user service: %v", err)
	}

	userJobProcessor, err := userJob.NewProcessor(l, cfg, cioClient, userService)
	if err != nil {
		log.Fatalf("failed to initialize user job processor: %v", err)
	}

	userHandler := worker.NewUserHTTPHandler(userJobProcessor)

	worker := NewWorker(userHandler)

	middlewares := []nethttp.StrictHTTPMiddlewareFunc{
		codegenkit.Recover,
	}

	handler := gen.HandlerFromMux(gen.NewStrictHandler(worker, middlewares), http.NewServeMux())
	if err := httpserverkit.Serve(handler, cfg.WorkerPort); err != nil {
		log.Fatal(err)
	}

	// FIXME: use hook to close all resources gracefully and in correct order
	cioClient.Close(ctx)
}
