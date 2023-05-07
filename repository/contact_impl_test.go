package repository

import (
	"contact-go/model"
	"testing"

	"github.com/stretchr/testify/suite"
)

type InMemoryRepoSuite struct {
	suite.Suite
	repo ContactRepository
}

func (s *InMemoryRepoSuite) SetupSuite() {
	model.Contacts = []model.Contact{
		{ID: 1, Name: "Reva", NoTelp: "555-1234-989"},
		{ID: 2, Name: "Tirta", NoTelp: "555-5678"},
		{ID: 3, Name: "Bagas", NoTelp: "555-9012"},
	}

	repo := NewContactRepository()
	s.repo = repo
}

func (s *InMemoryRepoSuite) TearDownSuite() {
	model.Contacts = []model.Contact{}
}

func TestInMemoryRepoSuite(t *testing.T) {
	suite.Run(t, new(InMemoryRepoSuite))
}

func (s *InMemoryRepoSuite) Test_contactRepository_List() {
	tests := []struct {
		name    string
		want    []model.Contact
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			want: []model.Contact{
				{ID: 2, Name: "Tirta", NoTelp: "555-5678"},
				{ID: 3, Name: "Bagas", NoTelp: "555-9012"},
				{ID: 4, Name: "Mixue", NoTelp: "555-9999"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.List()

			if s.Equal(tt.wantErr, err != nil, "contactRepository.List() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *InMemoryRepoSuite) Test_contactRepository_Add() {
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
					Name:   "Mixue",
					NoTelp: "555-9999",
				},
			},
			want: &model.Contact{
				ID:     4,
				Name:   "Mixue",
				NoTelp: "555-9999",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.Add(tt.args.newContact)

			if s.Equal(tt.wantErr, err != nil, "contactRepository.Add() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactRepository.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *InMemoryRepoSuite) Test_contactRepository_Detail() {
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
				Name:   "Tirta",
				NoTelp: "555-5678",
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

			if s.Equal(tt.wantErr, err != nil, "contactRepository.Detail() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactRepository.Detail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *InMemoryRepoSuite) Test_contactRepository_Update() {
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
					Name:   "Reva Iota",
					NoTelp: "555-1234-989",
				},
			},
			want: &model.Contact{
				ID:     3,
				Name:   "Reva Iota",
				NoTelp: "555-1234-989",
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

			if s.Equal(tt.wantErr, err != nil, "contactRepository.Update() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactRepository.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *InMemoryRepoSuite) Test_contactRepository_Delete() {
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				id: 1,
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

			s.Equal(tt.wantErr, err != nil, "contactRepository.Detail() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
