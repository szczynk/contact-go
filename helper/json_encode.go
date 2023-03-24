package helper

import (
	"contact-go/model"
	"encoding/json"
	"os"
)

func EncodeJSON(path string, contacts *[]model.Contact) error {
	writer, err := os.Create(path)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(writer)
	encoder.Encode(&contacts)
	return nil
}
