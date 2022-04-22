package implement

import (
	"github.com/tsukiamaoto/book-crawler-go/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (p *ProductRepository) GetProductByName(name string) (*model.Product, error) {
	var product *model.Product
	err := p.db.Model(&model.Product{}).Where("name = ?", name).Find(&product).Error

	return product, err
}

func (p *ProductRepository) AddProduct(product *model.Product) error {
	err := p.db.Create(product).Error
	if err != nil {
		return err
	}
	
	return nil
}