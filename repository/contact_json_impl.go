package repository

import (
	"contact-go/helper"
	"contact-go/model"
	"errors"
)

const jsonFile = "contact/contact.json"

type contactJsonRepository struct{}

func NewContactJsonRepository() ContactRepository {
	return new(contactJsonRepository)
}

func (repo *contactJsonRepository) List() ([]model.Contact, error) {
	err := helper.DecodeJSON(jsonFile, &model.Contacts)
	if err != nil {
		return []model.Contact{}, err
	}

	return model.Contacts, nil
}

func (repo *contactJsonRepository) getLastID() (int64, error) {
	contacts, err := repo.List()

	var tempID int64
	for _, v := range contacts {
		if tempID < v.ID {
			tempID = v.ID
		}
	}
	return tempID, err
}

func (repo *contactJsonRepository) getIndexByID(id int64) (int, error) {
	contacts, err := repo.List()
	if err != nil {
		return -1, err
	}

	for i, v := range contacts {
		if id == v.ID {
			return i, nil
		}
	}

	return -1, errors.New("ID tidak ditemukan")
}

func (repo *contactJsonRepository) Add(contact *model.Contact) (*model.Contact, error) {
	id, err := repo.getLastID()
	if err != nil {
		return &model.Contact{}, err
	}

	newContact := contact
	newContact.ID = id + 1

	model.Contacts = append(model.Contacts, *newContact)

	err = helper.EncodeJSON(jsonFile, &model.Contacts)
	if err != nil {
		return &model.Contact{}, err
	}

	return newContact, nil
}

func (repo *contactJsonRepository) Detail(id int64) (*model.Contact, error) {
	contacts, err := repo.List()
	if err != nil {
		return &model.Contact{}, err
	}

	index, err := repo.getIndexByID(id)
	if err != nil {
		return &model.Contact{}, err
	}

	contact := contacts[index]

	return &contact, nil
}

func (repo *contactJsonRepository) Update(id int64, contact *model.Contact) (*model.Contact, error) {
	contacts, err := repo.List()
	if err != nil {
		return &model.Contact{}, err
	}

	index, err := repo.getIndexByID(id)
	if err != nil {
		return &model.Contact{}, err
	}

	updatedContact := &contacts[index]
	updatedContact.Name = contact.Name
	updatedContact.NoTelp = contact.NoTelp

	err = helper.EncodeJSON(jsonFile, &model.Contacts)
	if err != nil {
		return &model.Contact{}, err
	}

	return updatedContact, nil
}

func (repo *contactJsonRepository) Delete(id int64) error {
	index, err := repo.getIndexByID(id)
	if err != nil {
		return err
	}

	model.Contacts = append(model.Contacts[:index], model.Contacts[index+1:]...)

	err = helper.EncodeJSON(jsonFile, &model.Contacts)
	if err != nil {
		return err
	}

	return nil
}
