package repository

import (
	"contact-go/model"
	"database/sql"
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormRepoSuite struct {
	suite.Suite
	gormDB  *gorm.DB
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    ContactRepository
}

func (s *GormRepoSuite) SetupTest() {
	var err error

	mockDB, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Require().NoError(err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	//* gorm.Config handle internally, which can not mock explisitly
	gormConf := new(gorm.Config)
	gormConf.PrepareStmt = true
	gormConf.Logger = logger.Default.LogMode(logger.Info)
	gormConf.SkipDefaultTransaction = true

	gormDb, err := gorm.Open(dialector, gormConf)
	if err != nil {
		s.Require().NoError(err)
	}

	sqlDB, err := gormDb.DB()
	if err != nil {
		s.Require().NoError(err)
	}

	repo := NewContactGormRepository(gormDb)

	s.gormDB = gormDb
	s.mockDB = sqlDB
	s.mockSQL = mockSQL
	s.repo = repo
}

func (s *GormRepoSuite) TearDownTest() {
	stmtManager, ok := s.gormDB.ConnPool.(*gorm.PreparedStmtDB)
	if ok {
		// close prepared statements for *current session*
		for _, stmt := range stmtManager.Stmts {
			log.Println(stmt)
			stmt.Close() // close the prepared statement
		}
	}

	defer s.mockDB.Close()
}

func TestGormRepoSuite(t *testing.T) {
	suite.Run(t, new(GormRepoSuite))
}

func (s *GormRepoSuite) Test_contactGormRepository_List() {
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

				s.ExpectPrepare(regexp.QuoteMeta(query)).
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
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(assert.AnError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed prepare statement",
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
			sqlQuery := `SELECT "id","name","no_telp" FROM "contacts" ORDER BY id ASC`

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			got, err := s.repo.List()
			log.Println("case:", tt.name, ", got:", got, ", error:", err)

			if s.Equal(tt.wantErr, err != nil, "contactGormRepository.List() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactGormRepository.List() = %v, want %v", got, tt.want)
			}

			if err := s.mockSQL.ExpectationsWereMet(); err != nil {
				s.Errorf(err, "there were unfulfilled expectations: %s")
			}
		})
	}
}

func (s *GormRepoSuite) Test_contactGormRepository_Add() {
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
				rows := s.NewRows([]string{"id"}).
					AddRow(int64(1))

				s.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().
					WithArgs("test", "555-555-3232").
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
			name: "invalid Name or NoTelp",
			args: args{
				contact: &model.Contact{
					Name:   "",
					NoTelp: "555-555-3232",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("", "555-555-3232").
					WillReturnError(assert.AnError)
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
			sqlQuery := `INSERT INTO "contacts" ("name","no_telp") VALUES ($1,$2) RETURNING "id"`

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			got, err := s.repo.Add(tt.args.contact)

			if s.Equal(tt.wantErr, err != nil, "contactGormRepository.Add() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactGormRepository.Add() = %v, want %v", got, tt.want)
			}

			if err := s.mockSQL.ExpectationsWereMet(); err != nil {
				s.Errorf(err, "there were unfulfilled expectations: %s")
			}
		})
	}
}

func (s *GormRepoSuite) Test_contactGormRepository_Detail() {
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
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(0)).
					WillReturnError(assert.AnError)
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
			sqlQuery := `SELECT "id","name","no_telp" FROM "contacts" WHERE "contacts"."id" = $1 ORDER BY "contacts"."id" LIMIT 1`

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			got, err := s.repo.Detail(tt.args.id)

			if s.Equal(tt.wantErr, err != nil, "contactGormRepository.Detail() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactGormRepository.Detail() = %v, want %v", got, tt.want)
			}

			if err := s.mockSQL.ExpectationsWereMet(); err != nil {
				s.Errorf(err, "there were unfulfilled expectations: %s")
			}
		})
	}
}

func (s *GormRepoSuite) Test_contactGormRepository_Update() {
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
				rows := s.NewRows([]string{"id", "name", "no_telp"}).
					AddRow(int64(1), "jangkrik", "555-555-4000")

				s.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().
					WithArgs("jangkrik", "555-555-4000", int64(1)).
					WillReturnRows(rows)
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
					Name:   "jangkrik",
					NoTelp: "555-555-4000",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("jangkrik", "555-555-4000", int64(0)).
					WillReturnError(assert.AnError)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed to prepare statement",
			args: args{
				id: 1,
				contact: &model.Contact{
					Name:   "jangkrik",
					NoTelp: "555-555-4000",
				},
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
			sqlQuery := `UPDATE "contacts" SET "name"=$1,"no_telp"=$2 WHERE id = $3 RETURNING "id","name","no_telp"`

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			got, err := s.repo.Update(tt.args.id, tt.args.contact)

			if s.Equal(tt.wantErr, err != nil, "contactGormRepository.Update() error = %v, wantErr %v", err, tt.wantErr) {
				s.Equal(tt.want, got, "contactGormRepository.Update() = %v, want %v", got, tt.want)
			}

			if err := s.mockSQL.ExpectationsWereMet(); err != nil {
				s.Errorf(err, "there were unfulfilled expectations: %s")
			}
		})
	}
}

func (s *GormRepoSuite) Test_contactGormRepository_Delete() {
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

				s.ExpectPrepare(regexp.QuoteMeta(query)).
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
				s.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(int64(0)).
					WillReturnError(assert.AnError)
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

				s.ExpectPrepare(regexp.QuoteMeta(query)).
					WillReturnError(err)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			sqlQuery := `DELETE FROM "contacts" WHERE "contacts"."id" = $1`

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			err := s.repo.Delete(tt.args.id)

			s.Equal(tt.wantErr, err != nil, "contactUsecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
