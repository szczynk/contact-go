package repository

import (
	"contact-go/model"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempJSONFile(t *testing.T, contacts []model.Contact) string {
	tempFile, err := os.CreateTemp("", "test_contact_*.json")
	assert.NoError(t, err)

	encoder := json.NewEncoder(tempFile)
	err = encoder.Encode(&contacts)
	assert.NoError(t, err)

	err = tempFile.Close()
	assert.NoError(t, err)

	return tempFile.Name()
}

func Test_contactJsonRepository_List(t *testing.T) {
	tests := []struct {
		name    string
		want    []model.Contact
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			want: []model.Contact{
				{ID: 1, Name: "Licht", NoTelp: "123456789"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contacts := []model.Contact{
				{ID: 1, Name: "Licht", NoTelp: "123456789"},
			}

			tempFile := createTempJSONFile(t, contacts)
			defer os.Remove(tempFile)

			repo := NewContactJsonRepository(tempFile)
			repo.(*contactJsonRepository).decodeJSON(tempFile, &model.Contacts)

			got, err := repo.List()

			if assert.Equal(t, tt.wantErr, err != nil, "contactJsonRepository.List() error = %v, wantErr %v", err, tt.wantErr) {
				assert.Equal(t, tt.want, got, "contactJsonRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contactJsonRepository_Add(t *testing.T) {
	type args struct {
		newContact *model.Contact
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Contact
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				newContact: &model.Contact{
					Name:   "Test1",
					NoTelp: "131-555-1",
				},
			},
			want: &model.Contact{
				ID:     2,
				Name:   "Test1",
				NoTelp: "131-555-1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contacts := []model.Contact{
				{ID: 1, Name: "Licht", NoTelp: "123456789"},
			}

			tempFile := createTempJSONFile(t, contacts)
			defer os.Remove(tempFile)

			repo := NewContactJsonRepository(tempFile)
			repo.(*contactJsonRepository).decodeJSON(tempFile, &model.Contacts)

			got, err := repo.Add(tt.args.newContact)

			if assert.Equal(t, tt.wantErr, err != nil, "contactJsonRepository.Add() error = %v, wantErr %v", err, tt.wantErr) {
				assert.Equal(t, tt.want, got, "contactJsonRepository.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contactJsonRepository_Detail(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Contact
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				id: 2,
			},
			want: &model.Contact{
				ID:     2,
				Name:   "Test1",
				NoTelp: "131-555-1",
			},
			wantErr: false,
		},
		{
			name: "failed",
			args: args{
				id: 3,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contacts := []model.Contact{
				{ID: 1, Name: "Licht", NoTelp: "123456789"},
				{ID: 2, Name: "Test1", NoTelp: "131-555-1"},
			}

			tempFile := createTempJSONFile(t, contacts)
			defer os.Remove(tempFile)

			repo := NewContactJsonRepository(tempFile)
			repo.(*contactJsonRepository).decodeJSON(tempFile, &model.Contacts)

			got, err := repo.Detail(tt.args.id)

			if assert.Equal(t, tt.wantErr, err != nil, "contactJsonRepository.Detail() error = %v, wantErr %v", err, tt.wantErr) {
				assert.Equal(t, tt.want, got, "contactJsonRepository.Detail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contactJsonRepository_Update(t *testing.T) {
	type args struct {
		id             int64
		updatedContact *model.Contact
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Contact
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				id: 2,
				updatedContact: &model.Contact{
					Name:   "Test1r",
					NoTelp: "131-555-1",
				},
			},
			want: &model.Contact{
				ID:     2,
				Name:   "Test1r",
				NoTelp: "131-555-1",
			},
			wantErr: false,
		},
		{
			name: "failed",
			args: args{
				id: 5,
				updatedContact: &model.Contact{
					Name:   "Test1ra",
					NoTelp: "131-555-1123",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contacts := []model.Contact{
				{ID: 1, Name: "Licht", NoTelp: "123456789"},
				{ID: 2, Name: "Test1", NoTelp: "131-555-1"},
			}

			tempFile := createTempJSONFile(t, contacts)
			defer os.Remove(tempFile)

			repo := NewContactJsonRepository(tempFile)
			repo.(*contactJsonRepository).decodeJSON(tempFile, &model.Contacts)

			got, err := repo.Update(tt.args.id, tt.args.updatedContact)

			if assert.Equal(t, tt.wantErr, err != nil, "contactJsonRepository.Update() error = %v, wantErr %v", err, tt.wantErr) {
				assert.Equal(t, tt.want, got, "contactJsonRepository.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contactJsonRepository_Delete(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Contact
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				id: 2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contacts := []model.Contact{
				{ID: 1, Name: "Licht", NoTelp: "123456789"},
				{ID: 2, Name: "Test1", NoTelp: "131-555-1"},
			}

			tempFile := createTempJSONFile(t, contacts)
			defer os.Remove(tempFile)

			repo := NewContactJsonRepository(tempFile)
			repo.(*contactJsonRepository).decodeJSON(tempFile, &model.Contacts)

			err := repo.Delete(tt.args.id)

			assert.Equal(t, tt.wantErr, err != nil, "contactJsonRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
