package order

type NewOrderRequest struct {
	ProductIDs []int  `json:"product_ids" validate:"required"`
	Phone      string `json:"phone,omitempty" validate:"required"`
}
