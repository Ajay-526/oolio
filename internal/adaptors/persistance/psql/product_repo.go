package psql

import (
	"context"
	"oolio/internal/core/model"
)

type ProductRepo struct {
	*Database
}

func NewProductRepo(db *Database) ProductRepo {
	return ProductRepo{db}
}

func (r *ProductRepo) Create(ctx context.Context, product *model.Product) error {

	query := `
	INSERT INTO products (id, name, price, category)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Price,
		product.Category,
	).Scan(&product.ID)
}

func (r *ProductRepo) GetByID(ctx context.Context, id string) (*model.Product, error) {

	query := `
	SELECT id, name, price, category
	FROM products
	WHERE id = $1
	`

	product := &model.Product{}

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Category,
		)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepo) GetAll(ctx context.Context) ([]model.Product, error) {

	query := `
	SELECT id, name, price, category
	FROM products
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product

	for rows.Next() {

		var p model.Product

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Category,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}
