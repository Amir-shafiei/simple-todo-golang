package db

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var once sync.Once
var DB *gorm.DB

func Connect() *gorm.DB {

	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			panic("Error loading .env file")
		}
		Username := os.Getenv("DB_USER")
		Password := os.Getenv("DB_PASSWORD")
		Host := os.Getenv("DB_HOST")
		Port := os.Getenv("DB_PORT")
		dsn := Username + ":" + Password + "@tcp(" + Host + ":" + Port + ")/test?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(dsn))

		if err != nil {
			panic(err)
		}
		DB = db

	})
	return DB

}
