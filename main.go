package main

import (
	"contact-go/handler"
	"contact-go/repository"
)

func main() {

	contactRepo := repository.NewContactRepository()
	ContactHandler := handler.NewcontactHandler(contactRepo)

	handler.Menu(ContactHandler)
}
