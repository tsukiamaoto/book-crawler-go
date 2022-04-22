package service

import (
	"github.com/tsukiamaoto/book-crawler-go/model"
	repo "github.com/tsukiamaoto/book-crawler-go/repository"
	service "github.com/tsukiamaoto/book-crawler-go/service/implement"
)

type Products interface {
	GetProductByName(name string) (*model.Product, error)
	AddProduct(product *model.Product) error
}

type Services struct {
	Products Products
}

type Repos struct {
	Repos *repo.Repositories
}

func NewServices(repos Repos) *Services {
	productsService := service.NewProductsService(repos.Repos.Products)

	return &Services{
		Products: productsService,
	}
}
