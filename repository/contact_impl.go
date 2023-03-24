package repository

import (
	"contact-go/model"
	"errors"
)

type contactRepository struct{}

func NewContactRepository() ContactRepository {
	return new(contactRepository)
}

func (repo *contactRepository) List() []model.Contact {
	return model.Contacts
}

func (repo *contactRepository) getLastID() int64 {
	contacts := repo.List()

	var tempID int64
	for _, v := range contacts {
		if tempID < v.ID {
			tempID = v.ID
		}
	}
	return tempID
}

func (repo *contactRepository) getIndexByID(id int64) (int, error) {
	contacts := repo.List()

	for i, v := range contacts {
		if id == v.ID {
			return i, nil
		}
	}

	return -1, errors.New("ID tidak ditemukan")
}

func (repo *contactRepository) Add(req model.ContactRequest) (model.Contact, error) {
	id := repo.getLastID()

	contact := model.Contact{
		ID:     id + 1,
		Name:   req.Name,
		NoTelp: req.NoTelp,
	}

	model.Contacts = append(model.Contacts, contact)

	return contact, nil
}

func (repo *contactRepository) Update(id int64, req model.ContactRequest) (model.Contact, error) {
	contacts := repo.List()
	index, err := repo.getIndexByID(id)

	if err != nil {
		return model.Contact{}, err
	}

	contact := &contacts[index]
	contact.Name = req.Name
	contact.NoTelp = req.NoTelp

	return *contact, nil
}

func (repo *contactRepository) Delete(id int64) error {
	index, err := repo.getIndexByID(id)

	if err != nil {
		return err
	}

	model.Contacts = append(model.Contacts[:index], model.Contacts[index+1:]...)

	return nil
}
