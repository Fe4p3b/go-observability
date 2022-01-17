package main

import (
	"encoding/json"
	"io"
)

type SignificantPerson struct {
	ID         string
	FirstName  string
	LastName   string
	Occupation string
}

func getPersonById(id string) *SignificantPerson {
	for _, v := range significantPeople {
		if v.ID == id {
			return v
		}
	}
	return nil
}

func (p *SignificantPerson) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

type SignificantPersons []*SignificantPerson

func getPersons() *SignificantPersons {
	return &significantPeople
}

func (p *SignificantPersons) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

var significantPeople = SignificantPersons{
	{"1", "Fyodor", "Dostoevsky", "Arts"},
	{"2", "Leo", "Tolstoy", "Arts"},
	{"3", "Jesus", "Christ", "Religion"},
	{"4", "Isaac", "Newton", "Science"},
	{"5", "Plato", "", "Philosophy"},
}
