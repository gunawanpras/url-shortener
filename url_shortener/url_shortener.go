package url_shortener

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gunawanpras/url-shortener/cache"
	"github.com/gunawanpras/url-shortener/config"
)

type URLService struct {
	store  cache.CacheImpl
	mutex  sync.Mutex
	config config.Config
}

func New(config config.Config, rCache cache.CacheImpl) *URLService {

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
		println("error", err.Error())
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
	mux.HandleFunc("/s", us.ShortenHandler)
	mux.HandleFunc("/s/", us.RedirectHandler)
	return mux
}

func (us *URLService) generateShortCode() (string, error) {
	return us.getRandomString(7)
}

func (us *URLService) getRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}
