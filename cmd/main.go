package main

import (
	"log"

	"github.com/platatest/internal/repository/factory"
	"github.com/platatest/internal/service"
	"github.com/platatest/pkg/config"
)

func main() {

	config, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config parsed")

	db, err := factory.NewPersistenceLayer(config.DatabaseURL(), factory.DBTYPE(config.Database.Type))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database connected")

	service.Serve(config, db)
}
