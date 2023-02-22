package models

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/project_assignment_synapsis_ecommerce"))
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&User{}, &Cart{}, &Order{})
	if err != nil {
		panic(err)
	}

	seedProduct(db)

	DB = db
}

func seedProduct(db *gorm.DB) {
	// seed data
	if err := db.AutoMigrate(&Product{}); err == nil && db.Migrator().HasTable(&Product{}) {
		if err := db.First(&Product{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			//Insert seed data
			db.Create(SeedData)
		}
	}
}
