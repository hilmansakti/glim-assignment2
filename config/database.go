package config

import (
	"assignment-2/models"
	"assignment-2/utils"
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbPort   = "5432"
	dbname   = "glim2"
	db       *gorm.DB
	store    *pgstore.PGStore
	err      error
)

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, dbPort)

	db, err = gorm.Open(postgres.Open(config), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("error connecting to database :", err)
	}
	// end connection
	db.AutoMigrate(&models.User{})

}

func GetDB() *gorm.DB {

	return db
}

func PgStore() *pgstore.PGStore {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, dbPort, dbname)
	authKey := utils.GenerateRandomKey(32)
	encryptionKey := utils.GenerateRandomKey(32)

	var err error
	store, err = pgstore.NewPGStore(url, authKey, encryptionKey)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(1)
	}

	return store
}
