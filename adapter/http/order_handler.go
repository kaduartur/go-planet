package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kaduartur/go-planet/adapter/http/request"
	"github.com/kaduartur/go-planet/core/port"
)

type orderHandler struct {
	order port.OrderService
}

func (h *orderHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	var createOrder request.CreateOrder
	if err := json.NewDecoder(r.Body).Decode(&createOrder); err != nil {
		responseString(w, http.StatusBadRequest, err.Error())
		return
	}

	cmd := createOrder.ToCommand()
	if err := h.order.Create(r.Context(), cmd); err != nil {
		responseString(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func makeOrderHandler(router chi.Router, orderService port.OrderService) {
	handler := &orderHandler{
		order: orderService,
	}

	// Order endpoints
	router.Route("/orders", func(r chi.Router) {
		r.Post("", handler.createOrder)
	})
}
