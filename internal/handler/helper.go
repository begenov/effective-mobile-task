package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) error(w http.ResponseWriter, code int, messageError string) {
	h.respond(w, code, map[string]string{"error": messageError})
}

func (h *Handler) respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
