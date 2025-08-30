package product

import (
	"api/orders/pkg/db"

	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	Database *db.Db
}

func NewProductkRepository(database *db.Db) *ProductRepository {
	return &ProductRepository{
		Database: database,
	}
}

func (repo *ProductRepository) Create(product *Product) (*Product, error) {
	result := repo.Database.DB.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (repo *ProductRepository) Update(product *Product) (*Product, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (repo *ProductRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ProductRepository) GetById(id uint) (*Product, error) {
	var product Product
	result := repo.Database.DB.First(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (repo *ProductRepository) List(page int, limit int) ([]Product, int64, error) {
	var products []Product
	var total int64

	offset := (page - 1) * limit

	if err := repo.Database.DB.Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := repo.Database.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&products)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return products, total, nil
}
