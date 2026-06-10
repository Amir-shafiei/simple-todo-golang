package db

import (
	"fmt"
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

		// چک ساده برای جلوگیری از crash بی‌دلیل
		if username == "" || password == "" || host == "" || port == "" || dbname == "" {
			panic("missing database environment variables")
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			username, password, host, port, dbname,
		)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		DB = db
	})

	return DB
}
