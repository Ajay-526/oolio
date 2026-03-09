package core

type PlaceOrderRequest struct {
	CouponCode string `json:"coupon_code"`
	Items      []Item `json:"items"`
}

type Item struct {
	ProductID string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
}

type PlaceOrderReponse struct {
	ID       string    `json:"id"`
	Items    []Item    `json:"items"`
	Products []Product `json:"products"`
}
