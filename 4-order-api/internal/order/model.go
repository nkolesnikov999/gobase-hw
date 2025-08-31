package order

import (
	"api/orders/internal/product"
	"api/orders/internal/user"

	"gorm.io/gorm"
)

// Order represents a customer's order belonging to a specific user and containing many products.
type Order struct {
	gorm.Model
	UserID   uint              `json:"user_id" gorm:"index"`
	User     user.User         `json:"-"`
	Products []product.Product `json:"products" gorm:"many2many:order_products;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
