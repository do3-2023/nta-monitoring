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
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Location    string `json:"location"`
}

var persons = []Person{
	{LastName: "Doe", PhoneNumber: "123456789", Location: "New York"},
	{LastName: "Smith", PhoneNumber: "987654321", Location: "Los Angeles"},
	{LastName: "Johnson", PhoneNumber: "456789123", Location: "Chicago"},
	{LastName: "Williams", PhoneNumber: "654321987", Location: "Houston"},
	{LastName: "Brown", PhoneNumber: "789123456", Location: "Phoenix"},
	{LastName: "Jones", PhoneNumber: "321987654", Location: "Philadelphia"},
	{LastName: "Garcia", PhoneNumber: "567891234", Location: "San Antonio"},
	{LastName: "Martinez", PhoneNumber: "891234567", Location: "San Diego"},
	{LastName: "Hernandez", PhoneNumber: "234567891", Location: "Dallas"},
}

func (db *DB) getPersons(w http.ResponseWriter, r *http.Request) {
	var persons []Person
	result := db.conn.Find(&persons)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(result.Error)
		return
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
