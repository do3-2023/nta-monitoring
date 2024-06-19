package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/fermyon/spin-go-sdk/pg"
	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
	"github.com/fermyon/spin/sdk/go/v2/variables"
)

type DB struct {
	DB *sql.DB
}

type Person struct {
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

var persons = []Person{
	{LastName: "Doe", PhoneNumber: "123456789"},
	{LastName: "Smith", PhoneNumber: "987654321"},
	{LastName: "Johnson", PhoneNumber: "456789123"},
	{LastName: "Brown", PhoneNumber: "654321987"},
	{LastName: "Williams", PhoneNumber: "789123456"},
	{LastName: "Jones", PhoneNumber: "321987654"},
	{LastName: "Garcia", PhoneNumber: "654123987"},
	{LastName: "Martinez", PhoneNumber: "987321654"},
	{LastName: "Hernandez", PhoneNumber: "123789456"},
	{LastName: "Gonzalez", PhoneNumber: "456321789"},
}

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		// create connection to DB
		url, _ := variables.Get("db_url")
		username, _ := variables.Get("db_username")
		password, _ := variables.Get("db_password")
		name, _ := variables.Get("db_name")
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
			url,
			username,
			password,
			name,
		)
		db, err := newDB(dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer db.DB.Close()
		fmt.Println("Connected to DB")

		rand.Seed(time.Now().UnixNano())

		router := spinhttp.NewRouter()
		router.GET("/healthz", db.checkDB)
		router.GET("/persons", db.getPersons)
		router.POST("/persons", db.addPerson)

		router.ServeHTTP(w, r)
	})
}

func main() {}

// Create a new DB connection
func newDB(dsn string) (*DB, error) {
	var db *sql.DB
	db = pg.Open(dsn)
	// create table for persons
	query := `
		CREATE TABLE IF NOT EXISTS persons (
			id SERIAL PRIMARY KEY,
			phone_number VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL
		)`
	_, err := db.Query(query)
	if err != nil {
		fmt.Println("Failed to create table:", err.Error())
		return nil, err
	}

	return &DB{db}, nil
}

// Check the DB connection by making a sql call
func (db *DB) checkDB(w http.ResponseWriter, r *http.Request, ps spinhttp.Params) {
	// Pinging the database
	err := db.DB.Ping()
	if err != nil {
		fmt.Println("Failed to ping database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Database is up")
}

func (db *DB) getPersons(w http.ResponseWriter, r *http.Request, ps spinhttp.Params) {
	err := db.DB.Ping()
	if err != nil {
		fmt.Println("Failed to ping database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("Getting persons")
	rows, err := db.DB.Query("SELECT last_name, phone_number FROM persons")
	if err != nil {
		fmt.Println(db.DB.Stats())
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var persons []*Person
	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.LastName, &person.PhoneNumber); err != nil {
			fmt.Println(err)
		}
		persons = append(persons, &person)
	}

	// convert to json
	out, err := json.Marshal(persons)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))
}

func (db *DB) addPerson(w http.ResponseWriter, r *http.Request, ps spinhttp.Params) {
	// get random person to add
	person := persons[rand.Intn(len(persons))]
	fmt.Println("Adding person", person)

	_, err := db.DB.Exec("INSERT INTO persons (last_name, phone_number) VALUES ($1, $2)", person.LastName, person.PhoneNumber)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// convert to json
	out, err := json.Marshal(person)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(out))
}
