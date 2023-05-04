package handler

import (
	"bytes"
	"contact-go/helper/logger"
	"contact-go/middleware"
	"contact-go/mocks"
	"contact-go/model"
	"contact-go/usecase"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func useMiddleware(handler http.HandlerFunc) *middleware.Middleware {
	l := logger.New(true)

	muxMiddleware := new(middleware.Middleware)
	muxMiddleware.Handler = handler

	muxMiddleware.Use(middleware.Cors)
	muxMiddleware.Use(middleware.ContentTypeJson)
	muxMiddleware.Use(
		func(w http.ResponseWriter, r *http.Request, next http.Handler) http.Handler {
			return middleware.Log(l, w, r, next)
		},
	)
	muxMiddleware.Use(
		func(w http.ResponseWriter, r *http.Request, next http.Handler) http.Handler {
			return middleware.Error(l, w, r, next)
		},
	)

	return muxMiddleware
}

func Test_contactHTTPHandler_List(t *testing.T) {
	tests := []struct {
		name       string
		UCResult   []model.Contact
		UCErr      error
		wantStatus int
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			UCResult: []model.Contact{
				{ID: 1, Name: "jaguar", NoTelp: "999-888-7777"},
				{ID: 2, Name: "Jane_Smith", NoTelp: "555-555-5678"},
				{ID: 3, Name: "jangkrik", NoTelp: "000-000-0000"},
			},
			UCErr:      nil,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContactUC := mocks.NewContactUsecase(t)
			mockContactUC.On("List").Return(tt.UCResult, tt.UCErr)

			h := NewContactHTTPHandler(mockContactUC)

			// Define your HTTP handler
			handler := http.HandlerFunc(h.List)
			method := "GET"
			url := "http://localhost:8080/contacts"

			m := useMiddleware(handler)

			req := httptest.NewRequest(method, url, nil)
			recorder := httptest.NewRecorder()

			m.ServeHTTP(recorder, req)

			response := recorder.Result()
			got := response.StatusCode

			if assert.Equal(t, tt.wantErr, tt.UCErr != nil, "ContactUsecase.List error = %v, wantErr %v", tt.UCErr, tt.wantErr) {
				assert.Equal(t, tt.wantStatus, got, "ContactHTTPHandler.List handler returned wrong status code: = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func Test_contactHTTPHandler_Add(t *testing.T) {
	type args struct {
		req *model.ContactRequest
	}
	tests := []struct {
		name       string
		args       args
		UCResult   *model.Contact
		UCErr      error
		wantStatus int
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				req: &model.ContactRequest{
					Name:   "test",
					NoTelp: "222-222-3232",
				},
			},
			UCResult: &model.Contact{
				ID:     1,
				Name:   "test",
				NoTelp: "222-222-3232",
			},
			UCErr:      nil,
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "invalid name",
			args: args{
				req: &model.ContactRequest{
					Name:   "",
					NoTelp: "222-222-3232",
				},
			},
			UCResult:   nil,
			UCErr:      assert.AnError,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid no_telp",
			args: args{
				req: &model.ContactRequest{
					Name:   "test",
					NoTelp: "",
				},
			},
			UCResult:   nil,
			UCErr:      assert.AnError,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContactRequest := new(model.ContactRequest)
			mockContactRequest.Name = tt.args.req.Name
			mockContactRequest.NoTelp = tt.args.req.NoTelp

			mockContactUC := mocks.NewContactUsecase(t)

			if !tt.wantErr {
				mockContactUC.On("Add", mockContactRequest).Return(tt.UCResult, tt.UCErr)
			}

			h := NewContactHTTPHandler(mockContactUC)

			// Define your HTTP handler
			handler := http.HandlerFunc(h.Add)
			method := "POST"
			url := "http://localhost:8080/contacts"

			m := useMiddleware(handler)

			reqBody, _ := json.Marshal(mockContactRequest)
			req := httptest.NewRequest(method, url, bytes.NewReader(reqBody))
			recorder := httptest.NewRecorder()

			m.ServeHTTP(recorder, req)

			response := recorder.Result()
			got := response.StatusCode

			log.Println(got, tt.wantStatus, tt.UCErr != nil, tt.wantErr)

			if assert.Equal(t, tt.wantErr, tt.UCErr != nil, "ContactUsecase.Add error = %v, wantErr %v", tt.UCErr, tt.wantErr) {
				assert.Equal(t, tt.wantStatus, got, "ContactHTTPHandler.Add handler returned wrong status code: = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func Test_contactHTTPHandler_Detail(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name       string
		args       args
		UCResult   *model.Contact
		UCErr      error
		wantStatus int
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				id: 1,
			},
			UCResult: &model.Contact{
				ID:     1,
				Name:   "test",
				NoTelp: "222-222-3232",
			},
			UCErr:      nil,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid ID",
			args: args{
				id: 0,
			},
			UCResult:   nil,
			UCErr:      sql.ErrNoRows,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContactUC := mocks.NewContactUsecase(t)

			if !tt.wantErr {
				mockContactUC.On("Detail", tt.args.id).Return(tt.UCResult, tt.UCErr)
			}

			h := NewContactHTTPHandler(mockContactUC)

			// Define your HTTP handler
			handler := http.HandlerFunc(h.Detail)
			method := "GET"
			url := fmt.Sprintf("http://localhost:8080/contacts/%v", tt.args.id)

			m := useMiddleware(handler)

			req := httptest.NewRequest(method, url, nil)
			recorder := httptest.NewRecorder()

			m.ServeHTTP(recorder, req)

			response := recorder.Result()
			got := response.StatusCode

			log.Println(got, tt.wantStatus, tt.UCErr != nil, tt.wantErr)

			if assert.Equal(t, tt.wantErr, tt.UCErr != nil, "ContactUsecase.Detail error = %v, wantErr %v", tt.UCErr, tt.wantErr) {
				assert.Equal(t, tt.wantStatus, got, "ContactHTTPHandler.Detail handler returned wrong status code: = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func Test_contactHTTPHandler_Update(t *testing.T) {
	type args struct {
		id  int64
		req *model.ContactRequest
	}
	tests := []struct {
		name       string
		args       args
		UCResult   *model.Contact
		UCErr      error
		wantStatus int
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				id: 1,
				req: &model.ContactRequest{
					Name:   "test1",
					NoTelp: "222-222-3232",
				},
			},
			UCResult: &model.Contact{
				ID:     1,
				Name:   "test1",
				NoTelp: "222-222-3232",
			},
			UCErr:      nil,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid ID",
			args: args{
				id: 0,
				req: &model.ContactRequest{
					Name:   "test1",
					NoTelp: "222-222-3232",
				},
			},
			UCResult:   nil,
			UCErr:      sql.ErrNoRows,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid name",
			args: args{
				id: 1,
				req: &model.ContactRequest{
					Name:   "",
					NoTelp: "222-222-3232",
				},
			},
			UCResult:   nil,
			UCErr:      assert.AnError,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid no_telp",
			args: args{
				id: 1,
				req: &model.ContactRequest{
					Name:   "test",
					NoTelp: "",
				},
			},
			UCResult:   nil,
			UCErr:      assert.AnError,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContactRequest := new(model.ContactRequest)
			mockContactRequest.Name = tt.args.req.Name
			mockContactRequest.NoTelp = tt.args.req.NoTelp

			mockContactUC := mocks.NewContactUsecase(t)

			if !tt.wantErr {
				mockContactUC.On("Update", tt.args.id, mockContactRequest).Return(tt.UCResult, tt.UCErr)
			}

			h := NewContactHTTPHandler(mockContactUC)

			// Define your HTTP handler
			handler := http.HandlerFunc(h.Update)
			method := "PATCH"
			url := fmt.Sprintf("http://localhost:8080/contacts/%v", tt.args.id)

			m := useMiddleware(handler)

			reqBody, _ := json.Marshal(mockContactRequest)
			req := httptest.NewRequest(method, url, bytes.NewReader(reqBody))
			recorder := httptest.NewRecorder()

			m.ServeHTTP(recorder, req)

			response := recorder.Result()
			got := response.StatusCode

			log.Println(got, tt.wantStatus, tt.UCErr != nil, tt.wantErr)

			if assert.Equal(t, tt.wantErr, tt.UCErr != nil, "ContactUsecase.Update error = %v, wantErr %v", tt.UCErr, tt.wantErr) {
				assert.Equal(t, tt.wantStatus, got, "ContactHTTPHandler.Update handler returned wrong status code: = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func Test_contactHTTPHandler_Delete(t *testing.T) {
	type fields struct {
		ContactUC usecase.ContactUsecase
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name       string
		args       args
		UCErr      error
		wantStatus int
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				id: 1,
			},
			UCErr:      nil,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid ID",
			args: args{
				id: 0,
			},
			UCErr:      sql.ErrNoRows,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContactUC := mocks.NewContactUsecase(t)

			if !tt.wantErr {
				mockContactUC.On("Delete", tt.args.id).Return(tt.UCErr)
			}

			h := NewContactHTTPHandler(mockContactUC)

			// Define your HTTP handler
			handler := http.HandlerFunc(h.Delete)
			method := "GET"
			url := fmt.Sprintf("http://localhost:8080/contacts/%v", tt.args.id)

			m := useMiddleware(handler)

			req := httptest.NewRequest(method, url, nil)
			recorder := httptest.NewRecorder()

			m.ServeHTTP(recorder, req)

			response := recorder.Result()
			got := response.StatusCode

			log.Println(got, tt.wantStatus, tt.UCErr != nil, tt.wantErr)

			if assert.Equal(t, tt.wantErr, tt.UCErr != nil, "ContactUsecase.Delete error = %v, wantErr %v", tt.UCErr, tt.wantErr) {
				assert.Equal(t, tt.wantStatus, got, "ContactHTTPHandler.Delete handler returned wrong status code: = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}
