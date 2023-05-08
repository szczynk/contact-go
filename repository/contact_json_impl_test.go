package repository

import (
	"contact-go/model"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type JsonRepoSuite struct {
	suite.Suite
	repo ContactRepository
}

func mockJsonFile(contacts *[]model.Contact, dir string, pattern string) (tempFileName string, err error) {
	tempFile, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return "", err
	}

	encoder := json.NewEncoder(tempFile)
	err = encoder.Encode(&contacts)
	if err != nil {
		return "", err
	}

	err = tempFile.Close()
	if err != nil {
		return "", err
	}

	tempFileName = tempFile.Name()
	return tempFileName, nil
}

func (s *JsonRepoSuite) SetupSuite() {
	contacts := []model.Contact{
		{ID: 1, Name: "Reva", NoTelp: "555-1234-989"},
		{ID: 2, Name: "Tirta", NoTelp: "555-5678"},
		{ID: 3, Name: "Bagas", NoTelp: "555-9012"},
	}

	tempFileName, err := mockJsonFile(&contacts, "", "test_contact_*.json")
	s.NoError(err)

	repo := NewContactJsonRepository(tempFileName)

	s.repo = repo
}

func (s *JsonRepoSuite) TearDownSuite() {
	model.Contacts = []model.Contact{}
	os.Remove(s.repo.(*contactJsonRepository).jsonFile)
}

func TestJsonRepoSuite(t *testing.T) {
	suite.Run(t, new(JsonRepoSuite))
}

func (s *JsonRepoSuite) Test_contactJsonRepository_List() {
	tests := []struct {
		name    string
		want    []model.Contact
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			want: []model.Contact{
				{ID: 1, Name: "Reva", NoTelp: "555-1234-989"},
				{ID: 3, Name: "Bagas", NoTelp: "555-9012"},
				{ID: 4, Name: "Test1", NoTelp: "131-555-1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.List()

			if s.Equal(tt.wantErr, err != nil, "contactJsonRepository.List() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactJsonRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *JsonRepoSuite) Test_contactJsonRepository_Add() {
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
				ID:     4,
				Name:   "Test1",
				NoTelp: "131-555-1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			// err := func(ttname) error {
			// 	if ttname == "getLastID error" {
			// 		return errors.New("getLastID mocked error")
			// 	}
			// 	_, err := s.repo.(*contactJsonRepository).getLastID()
			// 	return err
			// }()

			got, err := s.repo.Add(tt.args.newContact)

			if s.Equal(tt.wantErr, err != nil, "contactJsonRepository.Add() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactJsonRepository.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *JsonRepoSuite) Test_contactJsonRepository_Detail() {
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
				id: 3,
			},
			want: &model.Contact{
				ID:     3,
				Name:   "Bagas",
				NoTelp: "555-9012",
			},
			wantErr: false,
		},
		{
			name: "failed",
			args: args{
				id: 5,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.Detail(tt.args.id)

			if s.Equal(tt.wantErr, err != nil, "contactJsonRepository.Detail() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactJsonRepository.Detail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *JsonRepoSuite) Test_contactJsonRepository_Update() {
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
				id: 3,
				updatedContact: &model.Contact{
					Name:   "Test1r",
					NoTelp: "131-555-1345",
				},
			},
			want: &model.Contact{
				ID:     3,
				Name:   "Test1r",
				NoTelp: "131-555-1345",
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
		s.Run(tt.name, func() {
			got, err := s.repo.Update(tt.args.id, tt.args.updatedContact)

			if s.Equal(tt.wantErr, err != nil, "contactJsonRepository.Update() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactJsonRepository.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *JsonRepoSuite) Test_contactJsonRepository_Delete() {
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
		{
			name: "failed",
			args: args{
				id: 10,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repo.Delete(tt.args.id)

			s.Equal(tt.wantErr, err != nil, "contactJsonRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}

func Test_contactJsonRepository_deencodeJSON(t *testing.T) {
	contacts := []model.Contact{
		{ID: 1, Name: "Reva", NoTelp: "555-1234-989"},
		{ID: 2, Name: "Tirta", NoTelp: "555-5678"},
		{ID: 3, Name: "Bagas", NoTelp: "555-9012"},
	}

	tests := []struct {
		name    string
		dir     string
		pattern string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "success",
			pattern: "test_contact_*.json",
			wantErr: false,
		},
		{
			name:    "fail",
			dir:     "tmp/data",
			pattern: "test_contact_*.json",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonFile, err := mockJsonFile(&contacts, tt.dir, tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("mockJsonFile error = %v, wantErr %v", err, tt.wantErr)
			}

			repo := &contactJsonRepository{
				jsonFile: jsonFile,
			}
			if err := repo.encodeJSON(); (err != nil) != tt.wantErr {
				t.Errorf("contactJsonRepository.encodeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := repo.decodeJSON(); (err != nil) != tt.wantErr {
				t.Errorf("contactJsonRepository.decodeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
