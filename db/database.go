package db

import (
	Config "github.com/tsukiamaoto/book-crawler-go/config"
	models "github.com/tsukiamaoto/book-crawler-go/model"

	"fmt"

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
	// connected to postgres db just to be able create db statement
	postgresDB, err := gorm.Open(postgres.Open(conf.Databases["default"].Source))
	if err != nil {
		fmt.Println("使用 gorm 連線 DB 發生錯誤，原因為", err)
		log.Error(err)
		return nil, hasCreatedDB
	}

	// created traget database and connect to target database
	dbExec := fmt.Sprintf("CREATE DATABASE %s;", conf.Databases["shopCart"].Name)
	err = postgresDB.Exec(dbExec).Error
	if err == nil {
		hasCreatedDB = true
		fmt.Printf("建立 %s 資料庫\n", conf.Databases["shopCart"].Name)
	} else {
		fmt.Printf("%s 資料庫已經建立，連線該資料庫\n", conf.Databases["shopCart"].Name)
	}

	conn, err := gorm.Open(postgres.Open(conf.Databases["shopCart"].Source))
	if err != nil {
		fmt.Println("使用 gorm 連線 DB 發生錯誤，原因為", err)
		log.Error(err)
		return nil, hasCreatedDB
	}

	return conn, hasCreatedDB
}

func AutoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(new(*models.Category)); err != nil {
		panic("Category migration的失敗原因是" + err.Error())
	}
	fmt.Println("category db migration 成功！")

	if err := db.AutoMigrate(new(*models.Product)); err != nil {
		panic("資料庫Product migration的失敗原因是" + err.Error())
	}
	fmt.Println("product db migration 成功！")
}
