package helper

import (
	"contact-go/model"
	"encoding/json"
	"os"
)

func DecodeJSON(path string, contacts *[]model.Contact) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(reader)
	decoder.Decode(&contacts)
	return nil
}
