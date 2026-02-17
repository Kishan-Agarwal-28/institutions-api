package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	lru "github.com/hashicorp/golang-lru/v2"
)

var cache *lru.Cache[string, string]

func main() {
	// Initialize LRU cache
	var err error
	cache, err = lru.New[string, string](1000)
	if err != nil {
		log.Fatal(err)
	}

	// Load constellation data
	loadGeoJSON()
	if len(constellationData.Features) == 0 {
		panic("CRITICAL ERROR: No constellation features found in JSON!")
	}

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.Compress(5))
	r.Use(httprate.Limit(
		100,
		1*time.Minute,
		httprate.WithKeyFuncs(httprate.KeyByIP),
	))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	}))

	// Routes
	r.Get("/", documentationHandler)
	r.Get("/api/generate-avatar", generateAvatarHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("🚀 Server running on port %s", port)
	http.ListenAndServe(":"+port, r)
}
