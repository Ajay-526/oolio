package core

type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Category string `json:"category"`
}

type ProductList struct {
	Products []Product `json:"products"`
}

type CreateProductRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type CreateProductResponse struct {
	ID string `json:"id"`
}
