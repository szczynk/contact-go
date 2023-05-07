package repository

import (
	"contact-go/helper/apperrors"
	"contact-go/model"
	"encoding/json"
	"os"
)

type contactJsonRepository struct {
	jsonFile string
}

func NewContactJsonRepository(jsonFilePath string) ContactRepository {
	repo := new(contactJsonRepository)
	repo.jsonFile = jsonFilePath
	return repo
}

func (repo *contactJsonRepository) encodeJSON() error {
	writer, err := os.Create(repo.jsonFile)
	if err != nil {
		return err
	}
	defer writer.Close()

	encoder := json.NewEncoder(writer)
	err = encoder.Encode(&model.Contacts)
	if err != nil {
		return err
	}
	return nil
}

func (repo *contactJsonRepository) decodeJSON() error {
	reader, err := os.Open(repo.jsonFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&model.Contacts)
	if err != nil {
		return err
	}
	return nil
}

func (repo *contactJsonRepository) List() ([]model.Contact, error) {
	err := repo.decodeJSON()
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

	return -1, apperrors.NewAppError(apperrors.ErrContactNotFound)
}

func (repo *contactJsonRepository) Add(contact *model.Contact) (*model.Contact, error) {
	id, err := repo.getLastID()
	if err != nil {
		return nil, err
	}

	newContact := contact
	newContact.ID = id + 1

	model.Contacts = append(model.Contacts, *newContact)

	err = repo.encodeJSON()
	if err != nil {
		return nil, err
	}

	return newContact, nil
}

func (repo *contactJsonRepository) Detail(id int64) (*model.Contact, error) {
	contacts, err := repo.List()
	if err != nil {
		return nil, err
	}

	index, err := repo.getIndexByID(id)
	if err != nil {
		return nil, err
	}

	contact := contacts[index]

	return &contact, nil
}

func (repo *contactJsonRepository) Update(id int64, contact *model.Contact) (*model.Contact, error) {
	contacts, err := repo.List()
	if err != nil {
		return nil, err
	}

	index, err := repo.getIndexByID(id)
	if err != nil {
		return nil, err
	}

	updatedContact := &contacts[index]
	updatedContact.Name = contact.Name
	updatedContact.NoTelp = contact.NoTelp

	err = repo.encodeJSON()
	if err != nil {
		return nil, err
	}

	return updatedContact, nil
}

func (repo *contactJsonRepository) Delete(id int64) error {
	index, err := repo.getIndexByID(id)
	if err != nil {
		return err
	}

	model.Contacts = append(model.Contacts[:index], model.Contacts[index+1:]...)

	err = repo.encodeJSON()
	if err != nil {
		return err
	}

	return nil
}
