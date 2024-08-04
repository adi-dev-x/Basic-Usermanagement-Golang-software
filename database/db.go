package database

import (
	"fmt"
	"os"
	"projectgo/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("failed to connect to env")
	}
}
func DBconnect() {
	LoadEnv()
	dsn := os.Getenv("DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	DB = db
	fmt.Println("the db is createddd    ....///////")
	DB.AutoMigrate(&model.UserModel{}, &model.AdminModel{}, &model.ProductModel{}, &model.Bid{}, &model.CrediterModel{}, &model.VendorModel{}, &model.Credit{}, &model.Loans{}, &model.Payment{})
}
