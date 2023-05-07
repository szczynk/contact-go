package repository

import (
	"contact-go/model"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type JsonRepoSuite struct {
	suite.Suite
	repo ContactRepository
}

func (s *JsonRepoSuite) SetupSuite() {
	contacts := []model.Contact{
		{ID: 1, Name: "Reva", NoTelp: "555-1234-989"},
		{ID: 2, Name: "Tirta", NoTelp: "555-5678"},
		{ID: 3, Name: "Bagas", NoTelp: "555-9012"},
	}

	tempFile, err := os.CreateTemp("", "test_contact_*.json")
	s.Require().NoError(err)

	encoder := json.NewEncoder(tempFile)
	err = encoder.Encode(&contacts)
	s.Require().NoError(err)

	err = tempFile.Close()
	s.Require().NoError(err)

	tempFileName := tempFile.Name()

	repo := new(contactJsonRepository)
	repo.jsonFile = tempFileName

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
			err := s.repo.(*contactJsonRepository).decodeJSON()
			log.Printf("\ndecodeJSON err: %v\n", err)
			s.Require().NoError(err)

			got, err := s.repo.List()
			log.Printf("\nrepo.List err: %v\n", err)

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

			err := s.repo.(*contactJsonRepository).decodeJSON()
			log.Printf("\ndecodeJSON err: %v\n", err)
			s.Require().NoError(err)

			got, err := s.repo.Add(tt.args.newContact)
			log.Printf("\nrepo.Add err: %v\n", err)

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
			log.Printf("\nrepo.Detail err: %v\n", err)

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
			err := s.repo.(*contactJsonRepository).encodeJSON()
			log.Printf("\nencodeJSON err: %v\n", err)
			s.Require().NoError(err)

			got, err := s.repo.Update(tt.args.id, tt.args.updatedContact)
			log.Printf("\nrepo.Update err: %v\n", err)

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
			err := s.repo.(*contactJsonRepository).encodeJSON()
			log.Printf("\nencodeJSON err: %v\n", err)
			s.Require().NoError(err)

			err = s.repo.Delete(tt.args.id)
			log.Printf("\nrepo.Delete err: %v\n", err)

			s.Equal(tt.wantErr, err != nil, "contactJsonRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
