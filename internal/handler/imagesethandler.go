package handler

import "net/http"

type ImageSetHandler struct {
}

func (ish *ImageSetHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get ImageSet"))
}

func (ish *ImageSetHandler) Post(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Post ImageSet"))
}

func (ish *ImageSetHandler) Put(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Put ImageSet"))
}

func (ish *ImageSetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete ImageSet"))
}
