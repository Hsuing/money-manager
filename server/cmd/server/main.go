package main

import (
	"log"
	"fmt"
	"os"

	"easyaccounts-server/internal/repo"
	"easyaccounts-server/internal/router"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
}

func main() {
	var cfg Config
	cfgData, err := os.ReadFile("configs/config.yaml")
	if err == nil {
		yaml.Unmarshal(cfgData, &cfg)
	}
	
	dbPath := cfg.Database.Path
	if dbPath == "" {
		dbPath = "easyaccounts.db" // fallback
	}

	// Initialize database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := repo.InitSchema(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	repo.SeedData(db)

	r := router.SetupRouter(db, dbPath)
	
	port := cfg.Server.Port
	if port == 0 {
		port = 8081 // fallback
	}
	addr := fmt.Sprintf(":%d", port)
	
	log.Printf("Starting server on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
