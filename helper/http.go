package helper

import "net/http"

type Request struct {
	Method  string
	Handler http.Handler
}

func (req Request) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != req.Method {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	req.Handler.ServeHTTP(w, r)
}
