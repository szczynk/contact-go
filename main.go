package main

import (
	"contact-go/handler"
	"contact-go/repository"
	"contact-go/usecase"
)

func main() {
	contactRepo := repository.NewContactRepository()
	contactUC := usecase.NewContactUsecase(contactRepo)
	contactHandler := handler.NewContactHandler(contactUC)

	handler.Menu(contactHandler)
}
