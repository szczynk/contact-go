package repository

import (
	"contact-go/config/db"
	"contact-go/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type contactGormRepository struct {
	db *gorm.DB
}

func NewContactGormRepository(db *gorm.DB) ContactRepository {
	r := new(contactGormRepository)
	r.db = db

	return r
}

// db.ExecContext(...) function is used for executing SQL statements
// that do not return any rows, such as INSERT, UPDATE, and DELETE statements.

// On the other hand, the db.QueryRowContext(...) function is used for
// executing SQL queries that return a single row of result set.

func (repo *contactGormRepository) List() ([]model.Contact, error) {
	var contacts []model.Contact
	var err error

	ctx, cancel := db.NewContext()
	defer cancel()

	result := repo.db.WithContext(ctx).Select("id", "name", "no_telp").Order("id ASC").Find(&contacts)

	if err = result.Error; err != nil {
		return nil, err
	}

	return contacts, nil
}

func (repo *contactGormRepository) Add(contact *model.Contact) (*model.Contact, error) {
	ctx, cancel := db.NewContext()
	defer cancel()

	result := repo.db.WithContext(ctx).Select("Name", "NoTelp").Create(&contact)

	if err := result.Error; err != nil {
		return nil, err
	}

	return contact, nil
}

func (repo *contactGormRepository) Detail(id int64) (*model.Contact, error) {
	contact := new(model.Contact)

	ctx, cancel := db.NewContext()
	defer cancel()

	result := repo.db.WithContext(ctx).Select("ID", "Name", "NoTelp").First(&contact, id)

	if err := result.Error; err != nil {
		return nil, err
	}

	return contact, nil
}

func (repo *contactGormRepository) Update(id int64, contact *model.Contact) (*model.Contact, error) {
	ctx, cancel := db.NewContext()
	defer cancel()

	updatedContact := new(model.Contact)

	returning := clause.Returning{
		Columns: []clause.Column{{Name: "id"}, {Name: "name"}, {Name: "no_telp"}},
	}
	result := repo.db.WithContext(ctx).Model(&updatedContact).Clauses(returning).Where("id = ?", id).Updates(contact)

	if err := result.Error; err != nil {
		return nil, err
	}

	return updatedContact, nil
}

func (repo *contactGormRepository) Delete(id int64) error {
	contact := new(model.Contact)

	ctx, cancel := db.NewContext()
	defer cancel()

	result := repo.db.WithContext(ctx).Delete(&contact, id)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}
