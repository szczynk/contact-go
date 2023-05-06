package handler

import (
	"bytes"
	"contact-go/helper/input"
	"contact-go/mocks"
	"contact-go/model"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_contactHandler_List(t *testing.T) {
	tests := []struct {
		name     string
		UCResult []model.Contact
		UCErr    error
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			UCResult: []model.Contact{
				{ID: 1, Name: "jaguar", NoTelp: "999-888-7777"},
				{ID: 2, Name: "Jane_Smith", NoTelp: "555-555-5678"},
				{ID: 3, Name: "jangkrik", NoTelp: "000-000-0000"},
			},
			UCErr:   nil,
			wantErr: false,
		},
		{
			name:     "failed",
			UCResult: []model.Contact{},
			UCErr:    assert.AnError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader("")
			inputReader := input.NewInputReader(reader)

			mockContactUC := mocks.NewContactUsecase(t)

			mockContactUC.On("List").Return(tt.UCResult, tt.UCErr)

			h := NewContactHandler(mockContactUC, inputReader)

			h.List()

			assert.Equal(t, tt.wantErr, tt.UCErr != nil, "ContactUsecase.List error = %v, wantErr %v", tt.UCErr, tt.wantErr)
		})
	}
}

func captureStdout() (func(), chan string) {
	// Redirect output using an io.Pipe
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Capture output in a separate goroutine
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// Restore the original stdout
	return func() {
		w.Close()
		os.Stdout = oldStdout
	}, outC
}

func restoreStdout(restore func(), outC chan string) string {
	// Read the captured output
	restore()
	return <-outC
}

func Test_contactHandler_Add(t *testing.T) {
	type args struct {
		req *model.ContactRequest
	}
	tests := []struct {
		name     string
		args     args
		UCResult *model.Contact
		UCErr    error
		want     string
		wantErr  bool
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
			UCErr:   nil,
			want:    "Berhasil add contact with id",
			wantErr: false,
		},
		{
			name: "invalid name",
			args: args{
				req: &model.ContactRequest{
					Name:   "",
					NoTelp: "222-222-3232",
				},
			},
			UCResult: nil,
			UCErr:    assert.AnError,
			want:     "Name yang dimasukkan tidak valid",
			wantErr:  true,
		},
		{
			name: "invalid no_telp",
			args: args{
				req: &model.ContactRequest{
					Name:   "test",
					NoTelp: "",
				},
			},
			UCResult: nil,
			UCErr:    assert.AnError,
			want:     "NoTelp yang dimasukkan tidak valid",
			wantErr:  true,
		},
		{
			name: "invalid on usecase",
			args: args{
				req: &model.ContactRequest{
					Name:   "test",
					NoTelp: "1231451431",
				},
			},
			UCResult: nil,
			UCErr:    assert.AnError,
			want:     assert.AnError.Error(),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(fmt.Sprintf("%s\n%s\n", tt.args.req.Name, tt.args.req.NoTelp))
			inputReader := input.NewInputReader(reader)

			mockContactUC := mocks.NewContactUsecase(t)

			if !tt.wantErr || strings.Contains(tt.name, "usecase") {
				mockContactUC.On("Add", mock.Anything).Return(tt.UCResult, tt.UCErr)
			}

			h := NewContactHandler(mockContactUC, inputReader)

			restore, outC := captureStdout()
			h.Add()
			got := restoreStdout(restore, outC)

			if assert.Equal(t, tt.wantErr, tt.UCErr != nil, "ContactUsecase.Add error = %v, wantErr %v", tt.UCErr, tt.wantErr) {
				assert.Contains(t, got, tt.want, "Expected got to contain '%s', but got '%s'", tt.want, got)
			}
		})
	}
}

func Test_contactHandler_Detail(t *testing.T) {
	type args struct {
		idStr string
	}
	tests := []struct {
		name     string
		args     args
		UCResult *model.Contact
		UCErr    error
		want     string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				idStr: "1",
			},
			UCResult: &model.Contact{
				ID:     1,
				Name:   "test",
				NoTelp: "222-222-3232",
			},
			UCErr:   nil,
			want:    "ID : 		1\nNama : 		test\nNo.Telp : 	222-222-3232",
			wantErr: false,
		},
		{
			name: "invalid id",
			args: args{
				idStr: "0",
			},
			UCResult: nil,
			UCErr:    assert.AnError,
			want:     "ID yang dimasukkan tidak valid",
			wantErr:  true,
		},
		{
			name: "invalid on usecase",
			args: args{
				idStr: "1",
			},
			UCResult: nil,
			UCErr:    assert.AnError,
			want:     assert.AnError.Error(),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(fmt.Sprintf("%s\n", tt.args.idStr))
			inputReader := input.NewInputReader(reader)

			mockContactUC := mocks.NewContactUsecase(t)

			if !tt.wantErr || strings.Contains(tt.name, "usecase") {
				mockContactUC.On("Detail", mock.Anything).Return(tt.UCResult, tt.UCErr)
			}

			h := NewContactHandler(mockContactUC, inputReader)

			restore, outC := captureStdout()
			h.Detail()
			got := restoreStdout(restore, outC)

			if assert.Equal(t, tt.wantErr, tt.UCErr != nil, "ContactUsecase.Detail error = %v, wantErr %v", tt.UCErr, tt.wantErr) {
				assert.Contains(t, got, tt.want, "Expected got to contain '%s', but got '%s'", tt.want, got)
			}
		})
	}
}

func Test_contactHandler_Update(t *testing.T) {
	type args struct {
		idStr string
		req   *model.ContactRequest
	}
	tests := []struct {
		name     string
		args     args
		UCResult *model.Contact
		UCErr    error
		want     string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				idStr: "1",
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
			UCErr:   nil,
			want:    "Berhasil update contact with id",
			wantErr: false,
		},
		{
			name: "invalid id",
			args: args{
				idStr: "0",
				req: &model.ContactRequest{
					Name:   "test1",
					NoTelp: "222-222-3232",
				},
			},
			UCResult: nil,
			UCErr:    assert.AnError,
			want:     "ID yang dimasukkan tidak valid",
			wantErr:  true,
		},
		{
			name: "invalid name",
			args: args{
				idStr: "1",
				req: &model.ContactRequest{
					Name:   "",
					NoTelp: "222-222-3232",
				},
			},
			UCResult: nil,
			UCErr:    assert.AnError,
			want:     "Name yang dimasukkan tidak valid",
			wantErr:  true,
		},
		{
			name: "invalid no_telp",
			args: args{
				idStr: "1",
				req: &model.ContactRequest{
					Name:   "test",
					NoTelp: "",
				},
			},
			UCResult: nil,
			UCErr:    assert.AnError,
			want:     "NoTelp yang dimasukkan tidak valid",
			wantErr:  true,
		},
		{
			name: "invalid on usecase",
			args: args{
				idStr: "1",
				req: &model.ContactRequest{
					Name:   "test1",
					NoTelp: "222-222-3232",
				},
			},
			UCResult: nil,
			UCErr:    assert.AnError,
			want:     assert.AnError.Error(),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(fmt.Sprintf("%s\n%s\n%s\n", tt.args.idStr, tt.args.req.Name, tt.args.req.NoTelp))
			inputReader := input.NewInputReader(reader)

			mockContactUC := mocks.NewContactUsecase(t)

			if !tt.wantErr || strings.Contains(tt.name, "usecase") {
				mockContactUC.On("Update", mock.Anything, mock.Anything).Return(tt.UCResult, tt.UCErr)
			}

			h := NewContactHandler(mockContactUC, inputReader)

			restore, outC := captureStdout()
			h.Update()
			got := restoreStdout(restore, outC)

			if assert.Equal(t, tt.wantErr, tt.UCErr != nil, "ContactUsecase.Update error = %v, wantErr %v", tt.UCErr, tt.wantErr) {
				assert.Contains(t, got, tt.want, "Expected got to contain '%s', but got '%s'", tt.want, got)
			}
		})
	}
}

func Test_contactHandler_Delete(t *testing.T) {
	type args struct {
		idStr string
	}
	tests := []struct {
		name    string
		args    args
		UCErr   error
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				idStr: "1",
			},
			UCErr:   nil,
			want:    "Berhasil delete contact with id",
			wantErr: false,
		},
		{
			name: "invalid id",
			args: args{
				idStr: "0",
			},
			UCErr:   assert.AnError,
			want:    "ID yang dimasukkan tidak valid",
			wantErr: true,
		},
		{
			name: "invalid on usecase",
			args: args{
				idStr: "1",
			},
			UCErr:   assert.AnError,
			want:    assert.AnError.Error(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(fmt.Sprintf("%s\n", tt.args.idStr))
			inputReader := input.NewInputReader(reader)

			mockContactUC := mocks.NewContactUsecase(t)

			if !tt.wantErr || strings.Contains(tt.name, "usecase") {
				mockContactUC.On("Delete", mock.Anything).Return(tt.UCErr)
			}

			h := NewContactHandler(mockContactUC, inputReader)

			restore, outC := captureStdout()
			h.Delete()
			got := restoreStdout(restore, outC)

			if assert.Equal(t, tt.wantErr, tt.UCErr != nil, "ContactUsecase.Delete error = %v, wantErr %v", tt.UCErr, tt.wantErr) {
				assert.Contains(t, got, tt.want, "Expected got to contain '%s', but got '%s'", tt.want, got)
			}
		})
	}
}
