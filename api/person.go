package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type Person struct {
	LastName    string
	PhoneNumber string
	Location    string
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
	query := `SELECT last_name, phone_number, location FROM persons`
	rows, err := db.conn.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	var persons []Person
	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.LastName, &person.PhoneNumber, &person.Location); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		persons = append(persons, person)
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

	// add to db
	query := `INSERT INTO persons (last_name, phone_number, location) VALUES ($1, $2, $3)`
	_, err := db.conn.Exec(
		query,
		person.LastName,
		person.PhoneNumber,
		person.Location,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
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
