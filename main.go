package main

import (
	"log"

	"github.com/gadumitrachioaiei/slotserver/server"
)

func main() {
	s := server.New(":8080")
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
