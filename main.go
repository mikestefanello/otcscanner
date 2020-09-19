package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/mikestefanello/otcscanner/config"
	"github.com/mikestefanello/otcscanner/handlers"
	"github.com/mikestefanello/otcscanner/repository"
	"github.com/mikestefanello/otcscanner/router"
)

// TODO: Testing

func main() {
	// Load application configuration
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	// Create an order repository
	repo, err := repository.NewMongoOrderRepository(cfg.Mongo)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to repository: %s", err.Error()))
	}

	// Create an HTTP handler
	handler := handlers.NewHTTPHandler(cfg, repo)

	// Load the router
	r := router.NewRouter(cfg, handler)

	// Start the server
	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Hostname, cfg.HTTP.Port)
	log.Info().Str("on", addr).Msg("Server started")
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal().Err(err).Msg("Server terminated")
	}
}
