package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kaduartur/go-planet/core/port"
)

func MakeHandler(router *chi.Mux, orderService port.OrderService) {

	v1 := router.Route("/api/v1", func(r chi.Router) {
		return
	})
	makeOrderHandler(v1, orderService)
}

func responseString(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}
