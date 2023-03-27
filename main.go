package main

import (
	"contact-go/config"
	"contact-go/handler"
	"contact-go/repository"
	"contact-go/usecase"
	"log"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var contactRepo repository.ContactRepository

	switch config.Storage {
	case "json":
		contactRepo = repository.NewContactJsonRepository()
	default:
		contactRepo = repository.NewContactRepository()
	}

	contactUC := usecase.NewContactUsecase(contactRepo)
	contactHandler := handler.NewContactHandler(contactUC)

	handler.Menu(contactHandler)
}
