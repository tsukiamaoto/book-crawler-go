package db

import (
	Config "github.com/tsukiamaoto/book-crawler-go/config"
	model "github.com/tsukiamaoto/book-crawler-go/model"
	"github.com/tsukiamaoto/book-crawler-go/utils"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a type alias for gorm.DB
type DB = gorm.DB

func DbConnect() (*gorm.DB, bool) {
	hasCreatedDB := false
	conf := Config.LoadConfig()
	// // connected to postgres db just to be able create db statement
	// postgresDB, err := gorm.Open(postgres.Open(conf.Databases["default"].Source))
	// if err != nil {
	// 	fmt.Println("使用 gorm 連線 DB 發生錯誤，原因為", err)
	// 	log.Error(err)
	// 	return nil, hasCreatedDB
	// }

	// // created traget database and connect to target database
	// dbExec := fmt.Sprintf("CREATE DATABASE %s;", conf.Databases["shopCart"].Name)
	// err = postgresDB.Exec(dbExec).Error
	// if err == nil {
	// 	hasCreatedDB = true
	// 	fmt.Printf("建立 %s 資料庫\n", conf.Databases["shopCart"].Name)
	// } else {
	// 	fmt.Printf("%s 資料庫已經建立，連線該資料庫\n", conf.Databases["shopCart"].Name)
	// }

	conn, err := gorm.Open(postgres.Open(conf.Databases["shopCart"].Source))
	if err != nil {
		log.Error("使用 gorm 連線 DB 發生錯誤，原因為", err)
		return nil, hasCreatedDB
	}

	return conn, hasCreatedDB
}

func AutoMigrate(db *gorm.DB) {
	productMigration(db)
	migrateProductTypes2Types(db)
}

func productMigration(db *gorm.DB) {
	if err := db.AutoMigrate(new(*model.Product)); err != nil {
		log.Panic("資料庫Product migration的失敗原因是" + err.Error())
	}
	log.Println("product db migration 成功！")

	if err := db.AutoMigrate(new(*model.Category)); err != nil {
		log.Panic("Category migration的失敗原因是" + err.Error())
	}
	log.Println("category db migration 成功！")

	if err := db.AutoMigrate(new(*model.Type)); err != nil {
		log.Panic("type migration的失敗原因是" + err.Error())
	}
	log.Println("type db migration 成功！")
}

func migrateProductTypes2Types(db *gorm.DB) {
	products := make([]*model.Product, 0)

	if err := db.Model(&model.Product{}).Preload("Categories").Find(&products).Error; err != nil {
		log.Error("Failed to find products with Preload Categories, the reason is", err)
	}

	for _, product := range products {
		for _, category := range product.Categories {
			var keys []string
			relations := utils.RelationMap(category.Types)
			// append root key
			keys = append(keys, "root")
			keys = append(keys, category.Types...)
			types := utils.BuildTypes(keys, relations)

			var parentId *int
			var rootTypeId *int
			for index, typeValue := range types {
				var isExistedType model.Type
				if err := db.Model(&model.Type{}).Where("name = ?", typeValue.Name).Limit(1).Find(&isExistedType).Error; err != nil {
					log.Error("Failed to find type with name, the reason is ", err)
				}
				// if type not found, created a new type
				// else type has found, save id for next type value
				if (model.Type{}) == isExistedType {
					typeValue.ParentID = parentId
					if err := db.Model(&model.Type{}).Create(&typeValue).Error; err != nil {
						log.Error("Failed to create type, the reason is ", err)
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
			if rootTypeId != nil && category.TypeID == nil {
				category.TypeID = rootTypeId
				if err := db.Model(&model.Category{}).Where("id = ?", category.ID).Updates(category).Error; err != nil {
					log.Error("Faild to updated type of category, the reason is ", err)
				}
			}
		}
	}

	log.Println("Successfully migrate ProductTypes to Types!")
}
