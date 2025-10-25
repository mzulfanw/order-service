package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mzulfanw/order-service/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	order, err := h.service.CreateOrder(r.Context(), body.ProductID, body.Quantity)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"statusCode": http.StatusBadRequest,
			"error":      err.Error(),
			"message":    "Failed to create order",
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    order,
	})
}

func (h *OrderHandler) GetByProductID(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["productId"]

	orders, err := h.service.GetByProductID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    orders,
	})
}
