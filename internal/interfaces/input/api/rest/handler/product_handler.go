package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	core "oolio/internal/core/dto"
	"oolio/internal/usecase"
)

type ProductHandler struct {
	productHandler *usecase.ProductService
}

func NewProductHandler(productHandler *usecase.ProductService) *ProductHandler {
	return &ProductHandler{productHandler: productHandler}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req core.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := h.productHandler.CreateProduct(r.Context(), &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create product: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing product ID", http.StatusBadRequest)
		return
	}

	product, err := h.productHandler.GetProductByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get product: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productHandler.GetAllProducts(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get products: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
