package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	GetJokes() ([]*Joke, error)
	GetJokeById(int) (*Joke, error)
	// GetJokeRandom() (*Joke, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables file")
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createJokeTable()
}

func (s *PostgresStore) createJokeTable() error {
	query := `create table if not exists joke (
		id serial primary key,
		type text,
		setup text,
		punchline text,
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) GetJokeById(id int) (*Joke, error) {
	rows, err := s.db.Query("select * from joke where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoJoke(rows)
	}

	return nil, fmt.Errorf("joke %d not found", id)
}

func (s *PostgresStore) GetJokes() ([]*Joke, error) {
	rows, err := s.db.Query("select * from joke")
	if err != nil {
		return nil, err
	}

	jokes := []*Joke{}
	for rows.Next() {
		joke, err := scanIntoJoke(rows)
		if err != nil {
			return nil, err
		}
		jokes = append(jokes, joke)
	}

	return jokes, nil
}

func scanIntoJoke(rows *sql.Rows) (*Joke, error) {
	joke := new(Joke)
	err := rows.Scan(
		&joke.Id,
		&joke.Type,
		&joke.Setup,
		&joke.Punchline,
	)

	return joke, err
}
