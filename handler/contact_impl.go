package handler

import (
	"contact-go/helper"
	"contact-go/helper/input"
	"contact-go/model"
	"contact-go/usecase"
	"fmt"
	"strconv"
	"strings"
)

type contactHandler struct {
	ContactUC usecase.ContactUsecase
	Input     *input.InputReader
}

func NewContactHandler(contactUC usecase.ContactUsecase, input *input.InputReader) ContactHandler {
	contactHandler := new(contactHandler)
	contactHandler.ContactUC = contactUC
	contactHandler.Input = input

	return contactHandler
}

func (handler *contactHandler) List() {
	_ = helper.ClearTerminal()

	contacts, err := handler.ContactUC.List()

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("|---------------|-----------------------|-----------------------|\n")
		fmt.Printf("| ID\t\t| Nama\t\t\t| No.Telp\t\t|\n")
		fmt.Printf("|---------------|-----------------------|-----------------------|\n")

		for _, v := range contacts {
			fmt.Printf("| %d\t\t| %s\t\t| %s\t\t|\n", v.ID, v.Name, v.NoTelp)
		}
		fmt.Printf("|---------------|-----------------------|-----------------------|\n")
	}
}

func (handler *contactHandler) Add() {
	_ = helper.ClearTerminal()

	fmt.Print("Name = ")
	name, err := handler.Input.Scan()
	if err != nil || name == "" {
		fmt.Println("Name yang dimasukkan tidak valid")
		return
	}

	fmt.Print("NoTelp = ")
	noTelp, err := handler.Input.Scan()
	if err != nil || noTelp == "" {
		fmt.Println("NoTelp yang dimasukkan tidak valid")
		return
	}

	contactRequest := model.ContactRequest{
		Name:   name,
		NoTelp: noTelp,
	}

	contact, err := handler.ContactUC.Add(&contactRequest)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Berhasil add contact with id", contact.ID)
	}
}

func (handler *contactHandler) Detail() {
	_ = helper.ClearTerminal()

	fmt.Print("Contact ID = ")
	idStr, err := handler.Input.Scan()
	if err != nil {
		fmt.Println("ID yang dimasukkan tidak valid")
		return
	}

	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		fmt.Println("ID yang dimasukkan tidak valid")
		return
	}

	contact, err := handler.ContactUC.Detail(id)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("ID : \t\t%d\nNama : \t\t%s\nNo.Telp : \t%s\n", contact.ID, contact.Name, contact.NoTelp)
	}
}

func (handler *contactHandler) Update() {
	_ = helper.ClearTerminal()

	fmt.Print("ID = ")
	idStr, err := handler.Input.Scan()
	if err != nil {
		fmt.Println("ID yang dimasukkan tidak valid")
		return
	}

	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		fmt.Println("ID yang dimasukkan tidak valid")
		return
	}

	fmt.Print("Name = ")
	name, err := handler.Input.Scan()
	if err != nil || name == "" {
		fmt.Println("Name yang dimasukkan tidak valid")
		return
	}

	fmt.Print("NoTelp = ")
	noTelp, err := handler.Input.Scan()
	if err != nil || noTelp == "" {
		fmt.Println("NoTelp yang dimasukkan tidak valid")
		return
	}

	contactRequest := model.ContactRequest{
		Name:   name,
		NoTelp: noTelp,
	}

	contact, err := handler.ContactUC.Update(id, &contactRequest)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Berhasil update contact with id", contact.ID)
	}
}

func (handler *contactHandler) Delete() {
	_ = helper.ClearTerminal()

	fmt.Print("ID = ")
	idStr, err := handler.Input.Scan()
	if err != nil {
		fmt.Println("ID yang dimasukkan tidak valid")
		return
	}

	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		fmt.Println("ID yang dimasukkan tidak valid")
		return
	}

	err = handler.ContactUC.Delete(id)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Berhasil delete contact with id", id)
	}
}
