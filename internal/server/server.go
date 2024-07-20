package server

import (
	"net/http"
	"time"

	"github.com/pokemonpower92/collagegenerator/internal/handler"
)

type ImageSetServer struct {
	handler handler.Handler
	server  *http.Server
	mux     *http.ServeMux
}

func NewImageSetServer(h handler.Handler) *ImageSetServer {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &ImageSetServer{
		handler: h,
		server:  server,
		mux:     mux,
	}
}

func (iss *ImageSetServer) Start() {
	iss.mux.HandleFunc("GET /imagesets", iss.handler.Get)
	iss.mux.HandleFunc("POST /imagesets", iss.handler.Post)
	iss.mux.HandleFunc("PUT /imagesets", iss.handler.Put)
	iss.mux.HandleFunc("DELETE /imagesets", iss.handler.Delete)

	err := iss.server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
