package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DB interface {
	GetJokes(page int) ([]*Joke, error)
	GetJokeById(int) (*Joke, error)
	// GetJokeRandom() (*Joke, error)
}

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB() (*PostgresDB, error) {
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

	return &PostgresDB{
		db: db,
	}, nil
}

func (s *PostgresDB) Init() error {
	return s.createJokeTable()
}

func (s *PostgresDB) createJokeTable() error {
	query := `create table if not exists joke (
		id serial primary key,
		type text,
		setup text,
		punchline text,
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresDB) GetJokeById(id int) (*Joke, error) {
	rows, err := s.db.Query("select * from joke where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoJoke(rows)
	}

	return nil, fmt.Errorf("joke %d not found", id)
}

func (s *PostgresDB) GetJokes(page int) ([]*Joke, error) {
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * 20

	rows, err := s.db.Query(
		"SELECT * FROM joke LIMIT $1 OFFSET $2",
		20, offset)
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
