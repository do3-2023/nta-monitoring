package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	conn *gorm.DB
}

func main() {
	// create connection to DB
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Paris",
		os.Getenv("DB_URL"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := newDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to DB")

	// Migrate the schema
	db.conn.AutoMigrate(&Person{})

	db.conn.Migrator().DropColumn(&Person{}, "location")

	log.Println("Table persons is ready")

	// create api router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	// health route
	r.Get("/healthz", db.checkDB)

	r.Get("/persons", db.getPersons)
	r.Post("/persons", db.addPerson)

	log.Println("Starting api on port 3000")
	http.ListenAndServe(":3000", r)
}

// Create a new DB connection
func newDB(dsn string) (*DB, error) {
	// Number of retry attempts
	maxRetries := 10

	var db *gorm.DB
	var err error

	// Retry connecting to the database
	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Successfully connected
			fmt.Println("Connected to database on attempt", i)
			break
		}

		// Log the error and wait before retrying
		fmt.Printf("Attempt %d: Failed to connect to database: %v\n", i, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		fmt.Println("Exceeded maximum retry attempts. Exiting...")
		return nil, err
	}
	return &DB{db}, err
}

// Check the DB connection by making a sql call
func (db *DB) checkDB(w http.ResponseWriter, r *http.Request) {
	// Pinging the database
	sqlDB, err := db.conn.DB()
	if err != nil {
		fmt.Println("Failed to get database instance:", err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		fmt.Println("Failed to ping database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Database is up")
}
