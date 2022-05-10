package implement

import (
	"github.com/tsukiamaoto/book-crawler-go/model"
	repo "github.com/tsukiamaoto/book-crawler-go/repository"
)

type ProductsService struct {
	repo repo.Products
}

func NewProductsService(repo repo.Products) *ProductsService {
	return &ProductsService{
		repo: repo,
	}
}

func (p *ProductsService) GetProductByName(name string) (*model.Product, error) {
	return p.repo.GetProductByName(name)
}

func (p *ProductsService) AddProduct(product *model.Product) (*model.Product, error) {
	return p.repo.AddProduct(product)
}

func (p *ProductsService) AddTypeByCategoryId(categoryId uint, productTypes []string) error {
	return p.repo.AddTypeByCategoryId(categoryId, productTypes)
}
