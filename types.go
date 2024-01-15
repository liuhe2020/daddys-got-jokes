package main

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
