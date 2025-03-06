//go:generate oapi-codegen --config=../../internal/api/codegen/spec.yaml ../../doc/api/api.yaml
//go:generate oapi-codegen --config=../../internal/api/codegen/types.yaml ../../doc/api/api.yaml
package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"

	embed "github.com/AddMile/backend"
	gen "github.com/AddMile/backend/internal/api/codegen"
	api "github.com/AddMile/backend/internal/api/handler"
	"github.com/AddMile/backend/internal/app/user"
	"github.com/AddMile/backend/internal/config"

	codegenkit "github.com/AddMile/backend/internal/kit/codegen"
	httpserverkit "github.com/AddMile/backend/internal/kit/httpserver"
	loggerkit "github.com/AddMile/backend/internal/kit/logger"
	pgkit "github.com/AddMile/backend/internal/kit/pg"
	pubsubkit "github.com/AddMile/backend/internal/kit/pubsub"
	ratelimitkit "github.com/AddMile/backend/internal/kit/ratelimit"
	validatorkit "github.com/AddMile/backend/internal/kit/validator"
)

type API struct {
	*api.UserHTTPHandler
}

func NewAPI(
	userHandler *api.UserHTTPHandler,
) API {
	return API{
		userHandler,
	}
}

func main() {
	ctx := context.Background()
	cfg := config.Load()

	l := loggerkit.NewLogger(os.Stderr, cfg.Debug)

	validator := validatorkit.New()

	db, err := pgkit.New(cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}

	publisher, err := pubsubkit.NewPublisher(ctx, cfg.GoogleProjectID)
	if err != nil {
		log.Fatal(err)
	}

	userStorage := user.NewStorage(db)
	userQueue := user.NewQueue(cfg, publisher)
	userService, err := user.NewService(l, validator, userStorage, userQueue)
	if err != nil {
		log.Fatalf("initializing user service: %v", err)
	}
	userHandler := api.NewUserHTTPHandler(userService)

	api := NewAPI(userHandler)

	mux := http.NewServeMux()

	swagger, err := fs.Sub(fs.FS(embed.EmbedAPIAssets), "doc/api")
	if err != nil {
		log.Fatalf("initializing swagger assets: %v", err)
	}
	fileServer := http.FileServer(http.FS(swagger))
	mux.Handle("/", fileServer)

	middlewares := []nethttp.StrictHTTPMiddlewareFunc{
		codegenkit.RateLimit(ratelimitkit.New()),
		codegenkit.Recover,
	}

	handler := gen.HandlerFromMux(gen.NewStrictHandler(api, middlewares), mux)

	if err := httpserverkit.Serve(httpserverkit.CORS(cfg.CORSOrigin)(handler), cfg.APIPort); err != nil {
		log.Fatal(err)
	}
}
