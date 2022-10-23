package config

import (
	"fmt"
	"live-chat-gorilla/entity"
	"log"
	"net/http"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDatabaseConnection() *gorm.DB {

	envs, err := godotenv.Read(".env")

    if err != nil {
        log.Fatal("Error loading .env file")
    }

	dbHost := envs["DB_HOST"]
	dbPort := envs["DB_PORT"]
	dbUser := envs["DB_USER"]
	dbName := envs["DB_NAME"]
	
	// dbHost := "localhost"
	// dbPort := "3306"
	// dbName := "live_chat"
	// dbUser := "root"

	fmt.Println(dbHost)

	dsn := fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8&loc=Local", dbUser, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to create a connection to database")
	}
	entities := []interface{}{
		&entity.Presence{},
		&entity.Comment{},
	}
	errMigrate := db.AutoMigrate(entities...)
	if errMigrate != nil {
		panic("Failed to auto migrate")
	}
	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("there is error db close")
	}
	dbSQL.Close()
}

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Order("created_at desc").Offset(offset).Limit(pageSize)
	}
}
