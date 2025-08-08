package order

type NewOrderRequest struct {
	ProductIDs []uint `json:"product_ids" validate:"required,min=1,dive,required,min=1"`
}

type Response struct {
	ID         uint             `json:"id"`
	Phone      string           `json:"phone"`
	Status     string           `json:"status"`
	CreateAt   string           `json:"create_at"`
	Ordered    []ProductInOrder `json:"products"`
	TotalItems uint             `json:"total_items"`
}

type ProductInOrder struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
