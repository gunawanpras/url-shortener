package main

import (
	"log"
	"net/http"

	"github.com/gunawanpras/url-shortener/cache"
	"github.com/gunawanpras/url-shortener/config"
	"github.com/gunawanpras/url-shortener/url_shortener"
)

func main() {
	conf := config.LoadConfig("./config/config.yaml")
	log.Println("Load configuration...")

	redisCache := cache.NewRedisCache(conf.Cache)
	log.Println("Starting Redis Cache on port:", conf.Cache.Port)

	urlService := url_shortener.New(conf, redisCache)

	mux := urlService.Handler()
	http.Handle("/", mux)

	// Start HTTP server
	log.Println("Starting URL Shortener on port:", conf.Server.Port)
	if err := http.ListenAndServe(":"+conf.Server.Port, mux); err != nil {
		log.Fatal("Server error:", err)
	}
}
