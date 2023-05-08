package repository

import (
	"contact-go/model"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MysqlRepoSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    ContactRepository
}

func (s *MysqlRepoSuite) SetupTest() {
	var err error

	db, mock, err := sqlmock.New()
	if err != nil {
		s.Require().NoError(err)
	}

	repo := NewContactMysqlRepository(db)

	s.mockDB = db
	s.mockSQL = mock
	s.repo = repo
}

func (s *MysqlRepoSuite) TearDownTest() {
	defer s.mockDB.Close()
}

func TestMysqlRepoSuite(t *testing.T) {
	suite.Run(t, new(MysqlRepoSuite))
}

func (s *MysqlRepoSuite) Test_contactMysqlRepository_List() {
	tests := []struct {
		name       string
		beforeTest func(sqlmock.Sqlmock, string)
		want       []model.Contact
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				rows := s.NewRows([]string{"id", "name", "no_telp"}).
					AddRow(int64(1), "test", "555-555-3232")

				s.ExpectPrepare(query).
					ExpectQuery().
					WillReturnRows(rows)
			},
			want: []model.Contact{
				{ID: 1, Name: "test", NoTelp: "555-555-3232"},
			},
			wantErr: false,
		},
		{
			name: "failed",
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				s.ExpectPrepare(query).
					ExpectQuery().
					WillReturnError(assert.AnError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed prepare statement",
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				err := errors.New("prepare stmt error")

				s.ExpectPrepare(query).
					WillReturnError(err)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed rows scan",
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				rows := s.NewRows([]string{"id", "name", "no_telp"}).
					AddRow(int64(1), nil, "555-555-3232").
					RowError(1, errors.New("scanErr"))

				s.ExpectPrepare(query).
					ExpectQuery().
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed rows err",
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				rows := s.NewRows([]string{"id", "name", "no_telp"}).
					CloseError(errors.New("row error"))

				s.ExpectPrepare(query).
					ExpectQuery().
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			sqlQuery := "SELECT id, name, no_telp FROM contact ORDER BY id ASC"

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			got, err := s.repo.List()

			if s.Equal(tt.wantErr, err != nil, "contactMysqlRepository.List() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactMysqlRepository.List() = %v, want %v", got, tt.want)
			}

			if err := s.mockSQL.ExpectationsWereMet(); err != nil {
				s.Errorf(err, "there were unfulfilled expectations: %s")
			}
		})
	}
}

func (s *MysqlRepoSuite) Test_contactMysqlRepository_Add() {
	type args struct {
		contact *model.Contact
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock, string)
		want       *model.Contact
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				contact: &model.Contact{
					Name:   "test",
					NoTelp: "555-555-3232",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				stmt := s.ExpectPrepare(regexp.QuoteMeta(query))

				result := sqlmock.NewResult(1, 1)

				stmt.ExpectExec().
					WithArgs("test", "555-555-3232").
					WillReturnResult(result)
			},
			want: &model.Contact{
				ID:     1,
				Name:   "test",
				NoTelp: "555-555-3232",
			},
			wantErr: false,
		},
		{
			name: "invalid Name or NoTelp",
			args: args{
				contact: &model.Contact{
					Name:   "",
					NoTelp: "555-555-3232",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				stmt := s.ExpectPrepare(regexp.QuoteMeta(query))

				err := errors.New("invalid contact")

				stmt.ExpectExec().
					WithArgs("", "555-555-3232").
					WillReturnError(err)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed to prepare statement",
			args: args{
				contact: &model.Contact{
					Name:   "test",
					NoTelp: "555-555-3232",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				stmt := s.ExpectPrepare(regexp.QuoteMeta(query))

				err := errors.New("prepare stmt1 error")

				stmt.WillReturnError(err)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed to last insert id",
			args: args{
				contact: &model.Contact{
					Name:   "test",
					NoTelp: "555-555-3232",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				stmt := s.ExpectPrepare(regexp.QuoteMeta(query))

				result := sqlmock.NewErrorResult(errors.New("last insert id error"))

				stmt.ExpectExec().
					WithArgs("test", "555-555-3232").
					WillReturnResult(result)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			sqlQuery := "INSERT INTO contact(name, no_telp) VALUES (?, ?)"

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			got, err := s.repo.Add(tt.args.contact)

			if s.Equal(tt.wantErr, err != nil, "contactMysqlRepository.Add() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactMysqlRepository.Add() = %v, want %v", got, tt.want)
			}

			if err := s.mockSQL.ExpectationsWereMet(); err != nil {
				s.Errorf(err, "there were unfulfilled expectations: %s")
			}
		})
	}
}

func (s *MysqlRepoSuite) Test_contactMysqlRepository_Detail() {
	type args struct {
		id int64
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock, string)
		want       *model.Contact
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				id: 1,
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				rows := sqlmock.NewRows([]string{"id", "name", "no_telp"}).
					AddRow(int64(1), "test", "555-555-3232")

				s.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			want: &model.Contact{
				ID:     1,
				Name:   "test",
				NoTelp: "555-555-3232",
			},
			wantErr: false,
		},
		{
			name: "failed",
			args: args{
				id: 0,
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				s.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().
					WithArgs(int64(0)).
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed prepare statement",
			args: args{
				id: 0,
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				err := errors.New("prepare stmt error")

				s.ExpectPrepare(regexp.QuoteMeta(query)).
					WillReturnError(err)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			sqlQuery := "SELECT id, name, no_telp FROM contact WHERE id = ? LIMIT 1"

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			got, err := s.repo.Detail(tt.args.id)

			if s.Equal(tt.wantErr, err != nil, "contactMysqlRepository.Detail() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactMysqlRepository.Detail() = %v, want %v", got, tt.want)
			}

			if err := s.mockSQL.ExpectationsWereMet(); err != nil {
				s.Errorf(err, "there were unfulfilled expectations: %s")
			}
		})
	}
}

func (s *MysqlRepoSuite) Test_contactMysqlRepository_Update() {
	type args struct {
		id      int64
		contact *model.Contact
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock, string)
		want       *model.Contact
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				id: 1,
				contact: &model.Contact{
					Name:   "jangkrik",
					NoTelp: "555-555-4000",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				stmt := s.ExpectPrepare(regexp.QuoteMeta(query))

				result := sqlmock.NewResult(1, 1)

				stmt.ExpectExec().
					WithArgs("jangkrik", "555-555-4000", int64(1)).
					WillReturnResult(result)
			},
			want: &model.Contact{
				ID:     1,
				Name:   "jangkrik",
				NoTelp: "555-555-4000",
			},
			wantErr: false,
		},
		{
			name: "invalid id, name, or no_telp",
			args: args{
				id: 0,
				contact: &model.Contact{
					Name:   "test",
					NoTelp: "555-555-3232",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				stmt := s.ExpectPrepare(regexp.QuoteMeta(query))

				err := errors.New("invalid contact")

				stmt.ExpectExec().
					WithArgs("test", "555-555-3232", int64(0)).
					WillReturnError(err)

			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed to prepare statement",
			args: args{
				id: 1,
				contact: &model.Contact{
					Name:   "test",
					NoTelp: "555-555-3232",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				stmt := s.ExpectPrepare(regexp.QuoteMeta(query))

				err := errors.New("prepare stmt1 error")

				stmt.WillReturnError(err)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			sqlQuery := "UPDATE contact SET name = ?, no_telp = ? WHERE id = ?"

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			got, err := s.repo.Update(tt.args.id, tt.args.contact)

			if s.Equal(tt.wantErr, err != nil, "contactMysqlRepository.Update() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactMysqlRepository.Update() = %v, want %v", got, tt.want)
			}

			if err := s.mockSQL.ExpectationsWereMet(); err != nil {
				s.Errorf(err, "there were unfulfilled expectations: %s")
			}
		})
	}
}

func (s *MysqlRepoSuite) Test_contactMysqlRepository_Delete() {
	type args struct {
		id int64
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock, string)
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				id: 1,
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				result := sqlmock.NewResult(1, 1)

				s.ExpectPrepare(regexp.QuoteMeta("DELETE FROM contact WHERE id = ?")).
					ExpectExec().
					WithArgs(int64(1)).
					WillReturnResult(result)
			},
			wantErr: false,
		},
		{
			name: "failed",
			args: args{
				id: 0,
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				s.ExpectPrepare(regexp.QuoteMeta("DELETE FROM contact WHERE id = ?")).
					ExpectExec().
					WithArgs(int64(0)).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "failed prepare statement",
			args: args{
				id: 0,
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				err := errors.New("prepare stmt error")

				s.ExpectPrepare(regexp.QuoteMeta("DELETE FROM contact WHERE id = ?")).
					WillReturnError(err)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			sqlQuery := "DELETE FROM contact WHERE id = ?"

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			err := s.repo.Delete(tt.args.id)

			s.Equal(tt.wantErr, err != nil, "contactUsecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
