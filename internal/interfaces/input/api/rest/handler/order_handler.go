package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	core "oolio/internal/core/dto"
	"oolio/internal/usecase"
)

type OrderHandler struct {
	orderService *usecase.OrderService
}

func NewOrderHandler(orderService *usecase.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (o *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req core.PlaceOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := o.orderService.PlaceOrder(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to place order: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
