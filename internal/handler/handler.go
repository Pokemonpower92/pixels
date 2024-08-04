package handler

import "net/http"

type Handler interface {
	get(w http.ResponseWriter, r *http.Request)
	post(w http.ResponseWriter, r *http.Request)
	put(w http.ResponseWriter, r *http.Request)
	delete(w http.ResponseWriter, r *http.Request)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
