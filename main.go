package main

import (
	"log"
)

func main() {
	db, err := NewPostgresDB()
	if err != nil {
		log.Printf("error %s", err)
	}

	// if err := db.Init(); err != nil {
	// 	log.Printf("error %s", err)
	// }

	server := NewServer(":8080", db)
	server.Run()
}
