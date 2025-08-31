package order

type OrderCreateRequest struct {
	ProductIDs []uint `json:"product_ids" validate:"required,min=1,dive,gt=0"`
}

type OrderResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	CreatedAt string `json:"created_at"`
	Products  any    `json:"products"`
}

type OrderListResponse struct {
	Items []Order `json:"items"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
	Total int64   `json:"total"`
}
