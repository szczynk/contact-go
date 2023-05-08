package model

type Contact struct {
	ID     int64  `json:"id" gorm:"primarykey"`
	Name   string `json:"name"`
	NoTelp string `json:"no_telp"`
}

var Contacts []Contact

type ContactRequest struct {
	Name   string `json:"name"`
	NoTelp string `json:"no_telp"`
}
