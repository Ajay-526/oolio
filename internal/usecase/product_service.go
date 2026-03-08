package usecase

import (
	"context"
	"oolio/internal/adaptors/persistance/psql"
	core "oolio/internal/core/dto"
	"oolio/internal/core/model"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepo psql.ProductRepo
}

func NewProductService(
	productRepo psql.ProductRepo,
) ProductService {
	return ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *core.CreateProductRequest) (*core.CreateProductResponse, error) {
	id := uuid.New().String()
	newProduct := &model.Product{
		ID:       id,
		Name:     product.Name,
		Price:    product.Price,
		Category: product.Category,
	}
	if err := s.productRepo.Create(ctx, newProduct); err != nil {
		return nil, err
	}

	return &core.CreateProductResponse{
		ID: id,
	}, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id string) (*core.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &core.Product{
		ID:       product.ID,
		Name:     product.Name,
		Price:    int64(product.Price),
		Category: product.Category,
	}, nil
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]core.Product, error) {
	products, err := s.productRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []core.Product
	for _, p := range products {
		result = append(result, core.Product{
			ID:       p.ID,
			Name:     p.Name,
			Price:    int64(p.Price),
			Category: p.Category,
		})
	}

	return result, nil
}
