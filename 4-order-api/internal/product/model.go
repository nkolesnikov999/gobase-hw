package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
}

func NewProduct(name string, description string, images pq.StringArray) *Product {
	return &Product{
		Name:        name,
		Description: description,
		Images:      images,
	}
}
