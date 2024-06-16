package database

import (
	"database/sql"
	"fmt"
)

type Joke struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

type JokesResults struct {
	Total      int     `json:"total"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Results    []*Joke `json:"results"`
}

// func (s *service) createJokeTable() error {
// 	query := `CREATE TABLE IF NOT EXISTS joke (
// 		id serial primary key,
// 		type text,
// 		setup text,
// 		punchline text
// 	)`
// 	_, err := s.db.Exec(query)
// 	return err
// }

func (s *service) GetJokeById(id int) (*Joke, error) {
	rows, err := s.db.Query("SELECT * FROM joke WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoJoke(rows)
	}
	return nil, fmt.Errorf("joke %d not found", id)
}

func (s *service) GetJokes(page int) (*JokesResults, error) {
	offset := (page - 1) * 20
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM joke").Scan(&count)
	if err != nil {
		return nil, err
	}
	totalPages := (count + 20 - 1) / 20
	if page > totalPages {
		return nil, fmt.Errorf("requested page %d is greater than total pages %d", page, totalPages)
	}
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
	jokesResults := &JokesResults{
		Total:      count,
		TotalPages: totalPages,
		Page:       page,
		Results:    jokes,
	}
	return jokesResults, nil
}

func (s *service) GetJokeRandom() (*Joke, error) {
	rows, err := s.db.Query("SELECT * FROM joke ORDER BY RANDOM() LIMIT 1")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoJoke(rows)
	}
	return nil, fmt.Errorf("joke not found")
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
