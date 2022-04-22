package repository

import (
	"github.com/tsukiamaoto/book-crawler-go/model"
	repo "github.com/tsukiamaoto/book-crawler-go/repository/implement"

	"gorm.io/gorm"
)

type Products interface {
	GetProductByName(name string) (*model.Product, error)
	AddProduct(product *model.Product) error
}

type Repositories struct {
	Products Products
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Products: repo.NewProductRepository(db),
	}
}
