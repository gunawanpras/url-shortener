package url_shortener

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gunawanpras/url-shortener/cache"
	"github.com/gunawanpras/url-shortener/config"
	"github.com/gunawanpras/url-shortener/helper"
)

type URLService struct {
	store  cache.ICache
	mutex  sync.Mutex
	config config.Config
}

func New(config config.Config, rCache cache.ICache) *URLService {

	return &URLService{
		store:  rCache,
		config: config,
	}
}

func (us *URLService) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	us.mutex.Lock()
	shortCode, err := us.generateShortCode()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	err = us.store.SetValue(ctx, shortCode, request.URL, time.Duration(us.config.Cache.Ttl)*time.Minute)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	us.mutex.Unlock()

	response := map[string]string{
		"short_url": us.config.URLShortener.BaseURL + "/s/" + shortCode,
	}
	json.NewEncoder(w).Encode(response)
}

func (us *URLService) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[len("/s/"):]

	originalURL, err := us.store.GetValue(r.Context(), shortCode)
	if err != nil {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func (us *URLService) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/s", helper.Request{Method: http.MethodPost, Handler: http.HandlerFunc(us.ShortenHandler)})
	mux.Handle("/s/", helper.Request{Method: http.MethodGet, Handler: http.HandlerFunc(us.RedirectHandler)})
	return mux
}

func (us *URLService) generateShortCode() (string, error) {
	return helper.GetRandomString(7)
}
