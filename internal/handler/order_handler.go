package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/mzulfanw/order-service/internal/dto"
	"github.com/mzulfanw/order-service/internal/service"
	"github.com/mzulfanw/order-service/response"
)

type OrderHandler struct {
	service *service.OrderService
	logger  *logrus.Logger
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, ForceColors: true})
	logger.SetLevel(logrus.InfoLevel)
	return &OrderHandler{
		service: service,
		logger:  logger,
	}
}

var validate = validator.New()

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	h.logger.Infof("Incoming CreateOrder request: %s %s", r.Method, r.URL.Path)
	var body dto.CreateOrderDTO

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.logger.Warnf("Invalid request body: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(&body); err != nil {
		h.logger.Warnf("Validation error: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	order, err := h.service.CreateOrder(r.Context(), body.ProductID, body.Quantity)
	if err != nil {
		h.logger.Errorf("Failed to create order: %v", err)
		response.ErrorResponse(w, http.StatusBadRequest, "Failed to create order: "+err.Error())
		return
	}
	h.logger.Infof("Order created successfully: %+v", order)
	response.SuccessResponse(w, http.StatusCreated, order, "Order created successfully")
}

func (h *OrderHandler) GetByProductID(w http.ResponseWriter, r *http.Request) {
	h.logger.Infof("Incoming GetByProductID request: %s %s", r.Method, r.URL.Path)
	id, _ := mux.Vars(r)["productId"]
	h.logger.Infof("ProductID: %s", id)
	orders, err := h.service.GetByProductID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("Failed to get orders for product %s: %v", id, err)
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to get orders")
		return
	}
	h.logger.Infof("Orders retrieved successfully for product %s: %+v", id, orders)
	response.SuccessResponse(w, http.StatusOK, orders, "Orders retrieved successfully")
}
