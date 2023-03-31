package repository

import (
	"contact-go/config/db"
	"contact-go/model"
	"database/sql"
)

type contactMysqlRepository struct {
	db *sql.DB
}

func NewContactMysqlRepository(db *sql.DB) ContactRepository {
	return &contactMysqlRepository{
		db: db,
	}
}

func (repo *contactMysqlRepository) List() ([]model.Contact, error) {
	var contacts []model.Contact
	var contact model.Contact
	var err error

	ctx, cancel := db.NewMysqlContext()
	defer cancel()

	sqlQuery := "SELECT id, name, no_telp FROM contact ORDER BY id ASC"
	rows, err := repo.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return contacts, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&contact.ID, &contact.Name, &contact.NoTelp)
		if err != nil {
			return contacts, err
		}

		contacts = append(contacts, contact)
	}

	err = rows.Err()
	if err != nil {
		return contacts, err
	}

	return contacts, nil
}

func (repo *contactMysqlRepository) Add(contact *model.Contact) (*model.Contact, error) {
	newContact := new(model.Contact)
	var err error

	ctx, cancel := db.NewMysqlContext()
	defer cancel()

	sqlQuery := "INSERT INTO contact(name, no_telp) VALUES (?, ?)"
	row, err := repo.db.ExecContext(ctx, sqlQuery, contact.Name, contact.NoTelp)
	if err != nil {
		return newContact, err
	}

	contactID, err := row.LastInsertId()
	if err != nil {
		return newContact, err
	}

	newContact, err = repo.Detail(contactID)
	if err != nil {
		return newContact, err
	}

	return newContact, nil
}

func (repo *contactMysqlRepository) Detail(id int64) (*model.Contact, error) {
	var contact = new(model.Contact)
	var err error

	ctx, cancel := db.NewMysqlContext()
	defer cancel()

	sqlQuery := "SELECT id, name, no_telp FROM contact WHERE id = ? LIMIT 1"
	row := repo.db.QueryRowContext(ctx, sqlQuery, id)
	err = row.Scan(&contact.ID, &contact.Name, &contact.NoTelp)
	if err != nil {
		return contact, err
	}

	return contact, nil
}

func (repo *contactMysqlRepository) Update(id int64, contact *model.Contact) (*model.Contact, error) {
	updatedContact := new(model.Contact)
	var err error

	ctx, cancel := db.NewMysqlContext()
	defer cancel()

	sqlQuery := "UPDATE contact SET name = ?, no_telp = ? WHERE id = ?"
	row, err := repo.db.ExecContext(ctx, sqlQuery, contact.Name, contact.NoTelp, id)
	if err != nil {
		return updatedContact, err
	}

	contactID, err := row.LastInsertId()
	if err != nil {
		return updatedContact, err
	}

	updatedContact, err = repo.Detail(contactID)
	if err != nil {
		return updatedContact, err
	}

	return updatedContact, nil
}

func (repo *contactMysqlRepository) Delete(id int64) error {
	ctx, cancel := db.NewMysqlContext()
	defer cancel()

	sqlQuery := "DELETE FROM contact WHERE id = ?"
	_, err := repo.db.ExecContext(ctx, sqlQuery, id)
	if err != nil {
		return err
	}

	return nil
}
