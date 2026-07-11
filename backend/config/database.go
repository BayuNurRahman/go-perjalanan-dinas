package config

import (
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB adalah variabel global untuk akses database di repository nanti
var DB *gorm.DB

func ConnectDatabase() {
	LoadEnv()

	host := GetEnv("DB_HOST", "localhost")
	user := GetEnv("DB_USER", "postgres")
	password := GetEnv("DB_PASSWORD", "secret45")
	dbName := GetEnv("DB_NAME", "Traveldb")
	port := GetEnv("DB_PORT", "5434")
	sslMode := GetEnv("DB_SSLMODE", "disable")
	timeZone := GetEnv("DB_TIMEZONE", "Asia/Jakarta")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbName, port, sslMode, timeZone)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Gagal terhubung ke database", "error", err)
		os.Exit(1)
	}

	DB = database
	slog.Info("Koneksi database PostgreSQL berhasil terhubung", "host", host, "dbname", dbName)
}
