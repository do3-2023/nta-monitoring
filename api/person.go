package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	LastName    string  `json:"last_name"`
	PhoneNumber string  `json:"phone_number"`
	Location    *string `json:"location"`
}

var persons = []Person{
	{LastName: "Doe", PhoneNumber: "123456789", Location: NullableString("New York")},
	{LastName: "Smith", PhoneNumber: "987654321", Location: NullableString("Los Angeles")},
	{LastName: "Johnson", PhoneNumber: "456789123", Location: NullableString("Chicago")},
	{LastName: "Brown", PhoneNumber: "654321987", Location: NullableString("Houston")},
	{LastName: "Williams", PhoneNumber: "789123456", Location: NullableString("Phoenix")},
	{LastName: "Jones", PhoneNumber: "321987654", Location: NullableString("Philadelphia")},
	{LastName: "Garcia", PhoneNumber: "654123987", Location: NullableString("San Antonio")},
	{LastName: "Martinez", PhoneNumber: "987321654", Location: NullableString("San Diego")},
	{LastName: "Hernandez", PhoneNumber: "123789456", Location: NullableString("Dallas")},
	{LastName: "Gonzalez", PhoneNumber: "456321789", Location: NullableString("San Jose")},
}

func NullableString(x string) *string {
	return &x
}

func (db *DB) getPersons(w http.ResponseWriter, r *http.Request) {
	var persons []Person
	result := db.conn.Find(&persons)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(result.Error)
		return
	}

	// remove location field from persons
	for i := range persons {
		persons[i].Location = nil
	}

	// convert to json
	out, err := json.Marshal(persons)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))
}

func (db *DB) addPerson(w http.ResponseWriter, r *http.Request) {
	// get random person to add
	person := persons[rand.Intn(len(persons))]

	res := db.conn.Create(&person)
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(res.Error)
		return
	}
	log.Println("New person added")

	// convert to json
	out, err := json.Marshal(person)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(out))
}
