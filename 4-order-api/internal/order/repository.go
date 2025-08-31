package order

import (
	"api/orders/internal/product"
	"api/orders/pkg/db"
	"errors"

	"gorm.io/gorm"
)

type OrderRepository struct {
	Database *db.Db
}

func NewOrderRepository(database *db.Db) *OrderRepository {
	return &OrderRepository{Database: database}
}

// Create creates an order for a user and attaches the provided product IDs.
func (repo *OrderRepository) Create(userID uint, productIDs []uint) (*Order, error) {
	var products []product.Product
	if len(productIDs) > 0 {
		if err := repo.Database.DB.Find(&products, productIDs).Error; err != nil {
			return nil, err
		}
		if len(products) != len(productIDs) {
			return nil, errors.New("some products not found")
		}
	}

	ord := &Order{UserID: userID, Products: products}
	if err := repo.Database.DB.Create(ord).Error; err != nil {
		return nil, err
	}
	// Reload with products preloaded to return a complete object
	if err := repo.Database.DB.Preload("Products").First(ord, ord.ID).Error; err != nil {
		return nil, err
	}
	return ord, nil
}

// GetByID returns order by id with products preloaded
func (repo *OrderRepository) GetByID(id uint) (*Order, error) {
	var ord Order
	if err := repo.Database.DB.Preload("Products").First(&ord, id).Error; err != nil {
		return nil, err
	}
	return &ord, nil
}

// ListByUser returns paginated orders for a user with products preloaded
func (repo *OrderRepository) ListByUser(userID uint, page int, limit int) ([]Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	if err := repo.Database.DB.Model(&Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []Order
	query := repo.Database.DB.Preload("Products").Where("user_id = ?", userID).Limit(limit).Offset(offset).Order("id DESC").Find(&orders)
	if query.Error != nil && !errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return nil, 0, query.Error
	}
	return orders, total, nil
}
