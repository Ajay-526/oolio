package usecase

import (
	"context"
	"errors"
	"fmt"
	"oolio/internal/adaptors/persistance/psql"
	"oolio/internal/adaptors/persistance/redis"
	core "oolio/internal/core/dto"

	"github.com/google/uuid"
)

type OrderService struct {
	promoService   *PromoService
	redisService   *redis.RedisService
	productService *ProductService
	productRepo    *psql.ProductRepo
}

func NewOrderService(
	promoService *PromoService,
	redisService *redis.RedisService,
	productService *ProductService,
	productRepo *psql.ProductRepo,
) *OrderService {
	return &OrderService{
		promoService:   promoService,
		redisService:   redisService,
		productService: productService,
		productRepo:    productRepo,
	}
}

func (s *OrderService) PlaceOrder(
	ctx context.Context,
	req core.PlaceOrderRequest,
) (*core.PlaceOrderReponse, error) {

	if len(req.Items) == 0 {
		return nil, errors.New("items cannot be empty")
	}

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %s", item.ProductID)
		}
	}

	// Promo validation
	if req.CouponCode != "" {
		valid := s.promoService.ValidatePromo(ctx, req.CouponCode)

		if !valid {
			return nil, errors.New("invalid promo code")
		}
	}

	// Collect product ids
	productIDs := make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	// Fetch products

	var products []core.Product

	for _, id := range productIDs {
		product, err := s.productService.GetProductByID(ctx, id)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}

	if len(products) != len(productIDs) {
		return nil, errors.New("some products not found")
	}

	resp := &core.PlaceOrderReponse{
		ID:       generateOrderID(),
		Items:    req.Items,
		Products: products,
	}

	return resp, nil
}

func generateOrderID() string {
	return fmt.Sprintf("order-%s", uuid.New().String())
}
