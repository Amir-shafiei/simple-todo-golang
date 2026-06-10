package db

import (
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var once sync.Once
var DB *gorm.DB

func Connect() *gorm.DB {

	once.Do(func() {

		username := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		dbname := os.Getenv("DB_NAME")

		dsn := username + ":" + password +
			"@tcp(" + host + ":" + port + ")/" + dbname +
			"?charset=utf8mb4&parseTime=True&loc=Local"

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		DB = db
	})

	return DB
}
