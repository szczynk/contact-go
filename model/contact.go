package model

type Contact struct {
	ID     int64
	Name   string
	NoTelp string
}

var Contacts []Contact = []Contact{
	{
		ID:     1,
		Name:   "adli",
		NoTelp: "0984738",
	},
	{
		ID:     2,
		Name:   "abdul",
		NoTelp: "0984738",
	},
	{
		ID:     3,
		Name:   "rizki",
		NoTelp: "0984738",
	},
	{
		ID:     4,
		Name:   "siapa",
		NoTelp: "0984738",
	},
	{
		ID:     5,
		Name:   "rokuy",
		NoTelp: "0984738",
	},
	{
		ID:     6,
		Name:   "fian",
		NoTelp: "0984738",
	},
	{
		ID:     7,
		Name:   "nandy",
		NoTelp: "0984738",
	},
	{
		ID:     8,
		Name:   "bagas",
		NoTelp: "0984738",
	},
	{
		ID:     9,
		Name:   "arbi",
		NoTelp: "0984738",
	},
	{
		ID:     10,
		Name:   "adam",
		NoTelp: "0984738",
	},
}

type ContactRequest struct {
	Name   string
	NoTelp string
}
