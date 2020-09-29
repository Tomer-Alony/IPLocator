package main

import (
	"github.com/Tomer-Alony/IPLocator/src/config"
	"github.com/Tomer-Alony/IPLocator/src/handlers/ip"
	"github.com/Tomer-Alony/IPLocator/src/middlewares/request_limit"
	"github.com/Tomer-Alony/IPLocator/src/services"
	"github.com/Tomer-Alony/IPLocator/src/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

func runServer() {
	cfg := config.LoadConfig()

	concurrencyManager := request_limit.NewConcurrencyManager(cfg.RateLimit, cfg.MaxTokenDuration)

	// Connect to store
	ipDataStore := store.NewAPIDataStore().Access(cfg.StorePath, cfg.StoreKey)

	// Init repos
	ipRepo := services.NewIPServiceRepo(ipDataStore)

	// Init routes
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		request_limit.ConcurrencyHandler(concurrencyManager),
		render.SetContentType(render.ContentTypeJSON),
	)

	router.Route("/", func(r chi.Router){
		// Add all submodules handlers by their group
		r.Mount("/v1", ip.Routes(ipRepo))
	})

	log.Println("Starting IPLocator server at", cfg.Port)
	if err := http.ListenAndServe(":" + cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}

func main ()  {
	runServer()
}