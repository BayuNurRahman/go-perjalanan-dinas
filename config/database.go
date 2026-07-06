package config

import (
	"fmt"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB adalah variabel global untuk akses database di repository nanti
var DB *gorm.DB

func ConnectDatabase() {
	host := "localhost"
	user := "postgres"
	password := "secret45" 
	dbName := "Traveldb"
	port := "5434"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, port)
	
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	DB = database
	fmt.Println("Koneksi database PostgreSQL berhasil terhubung!")
}