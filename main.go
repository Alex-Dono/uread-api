package main

import (
	"fmt"
	"net/http"
	"os"
	"uread/config"
	"uread/model"
	"uread/router"

	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config.SetupConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Config.Database.DBAddr,
		config.Config.Database.DBUser,
		config.Config.Database.DBPass,
		config.Config.Database.DBName,
		config.Config.Database.DBPort,
		config.Config.Database.DBType,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = db.AutoMigrate(&model.Book{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to auto-migrate: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Database migration successful\n")

	mux := router.SetupRouter(db)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8081"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	}).Handler(mux)

	fmt.Printf("Connecting on port 8080\n")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to server: %v\n", err)
		os.Exit(1)
	}
}
