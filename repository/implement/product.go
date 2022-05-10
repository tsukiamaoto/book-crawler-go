package implement

import (
	"github.com/tsukiamaoto/book-crawler-go/model"
	"github.com/tsukiamaoto/book-crawler-go/utils"

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

func (p *ProductRepository) AddProduct(product *model.Product) (*model.Product, error) {
	err := p.db.Create(product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductRepository) AddTypeByCategoryId(categoryId uint, productTypes []string) error {
	var keys []string
	relations := utils.RelationMap(productTypes)
	// append root key
	keys = append(keys, "root")
	keys = append(keys, productTypes...)
	types := utils.BuildTypes(keys, relations)

	var parentId *int
	var rootTypeId *int
	for index, typeValue := range types {
		var isExistedType model.Type
		if err := p.db.Model(&model.Type{}).Where("name = ?", typeValue.Name).Limit(1).Find(&isExistedType).Error; err != nil {
			return err
		}
		// if type not found, created a new type
		// else type has found, save id for next type value
		if (model.Type{}) == isExistedType {
			typeValue.ParentID = parentId
			if err := p.db.Model(&model.Type{}).Create(&typeValue).Error; err != nil {
				return err
			}
			// save parent id for next type
			parentId = &typeValue.ID
		} else {
			parentId = &isExistedType.ID
		}

		// save root id for type id of category
		if index == 0 {
			rootTypeId = parentId
		}
	}

	// updated type of categroies
	if rootTypeId != nil {
		if err := p.db.Model(&model.Category{}).Where("id = ?", categoryId).Update("TypeID", rootTypeId).Error; err != nil {
			return err
		}
	}
	return nil
}
