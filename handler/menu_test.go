package handler

import (
	"contact-go/helper/input"
	"contact-go/mocks"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMenu_ShowMenu(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		method  string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "success list",
			method:  "List",
			input:   "1\n6",
			want:    "Contact list",
			wantErr: false,
		},
		{
			name:    "success add",
			method:  "Add",
			input:   "2\n6",
			want:    "Add a new contact",
			wantErr: false,
		},
		{
			name:    "success detail",
			method:  "Detail",
			input:   "3\n6",
			want:    "Contact detail",
			wantErr: false,
		},
		{
			name:    "success update",
			method:  "Update",
			input:   "4\n6",
			want:    "Update a contact",
			wantErr: false,
		},
		{
			name:    "success delete",
			method:  "Delete",
			input:   "5\n6",
			want:    "Delete a contact",
			wantErr: false,
		},
		{
			name:    "back to menu",
			method:  "List",
			input:   "\n6",
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			inputReader := input.NewInputReader(reader)

			mockContactHandler := mocks.NewContactHandler(t)
			if !strings.Contains(tt.name, "menu") {
				mockContactHandler.On(tt.method).Return()
			}

			clear := func() error { return nil }
			showMenuList := func() {}

			m := NewMenu(mockContactHandler, inputReader, clear, showMenuList)

			restore, outC := captureStdout()
			err := m.ShowMenu()

			got := restoreStdout(restore, outC)

			log.Println("tt.name:", tt.name)
			log.Println("err:", err)
			log.Println("got:", got)
			assert.Equal(t, tt.wantErr, err != nil, "ContactUsecase.Add error = %v, wantErr %v", err, tt.wantErr)
			assert.Contains(t, got, tt.want, "Expected got to contain '%s', but got '%s'", tt.want, got)
		})
	}
}
