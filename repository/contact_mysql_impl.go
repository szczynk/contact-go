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

// db.ExecContext(...) function is used for executing SQL statements
// that do not return any rows, such as INSERT, UPDATE, and DELETE statements.

// On the other hand, the db.QueryRowContext(...) function is used for
// executing SQL queries that return a single row of result set.

func (repo *contactMysqlRepository) List() ([]model.Contact, error) {
	var contacts []model.Contact
	var contact model.Contact
	var err error

	ctx, cancel := db.NewContext()
	defer cancel()

	sqlQuery := "SELECT id, name, no_telp FROM contact ORDER BY id ASC"
	stmt, err := repo.db.PrepareContext(ctx, sqlQuery)
	if err != nil {
		return contacts, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
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
	ctx, cancel := db.NewContext()
	defer cancel()

	sqlQuery1 := "INSERT INTO contact(name, no_telp) VALUES (?, ?)"
	stmt, err := repo.db.PrepareContext(ctx, sqlQuery1)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row, err := stmt.ExecContext(ctx, contact.Name, contact.NoTelp)
	if err != nil {
		return nil, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}

	newContact := new(model.Contact)
	newContact.ID = id
	newContact.Name = contact.Name
	newContact.NoTelp = contact.NoTelp

	return newContact, nil
}

func (repo *contactMysqlRepository) Detail(id int64) (*model.Contact, error) {
	contact := new(model.Contact)
	var err error

	ctx, cancel := db.NewContext()
	defer cancel()

	sqlQuery := "SELECT id, name, no_telp FROM contact WHERE id = ? LIMIT 1"
	stmt, err := repo.db.PrepareContext(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	err = row.Scan(&contact.ID, &contact.Name, &contact.NoTelp)
	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (repo *contactMysqlRepository) Update(id int64, contact *model.Contact) (*model.Contact, error) {
	ctx, cancel := db.NewContext()
	defer cancel()

	sqlQuery := "UPDATE contact SET name = ?, no_telp = ? WHERE id = ?"
	stmt, err := repo.db.PrepareContext(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, contact.Name, contact.NoTelp, id)
	if err != nil {
		return nil, err
	}

	updatedContact := new(model.Contact)
	updatedContact.ID = id
	updatedContact.Name = contact.Name
	updatedContact.NoTelp = contact.NoTelp

	return updatedContact, nil
}

func (repo *contactMysqlRepository) Delete(id int64) error {
	ctx, cancel := db.NewContext()
	defer cancel()

	sqlQuery := "DELETE FROM contact WHERE id = ?"
	stmt, err := repo.db.PrepareContext(ctx, sqlQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
