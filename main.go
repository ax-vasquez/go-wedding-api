package main

import (
	"log"

	"github.com/ax-vasquez/wedding-site-api/db"
	"github.com/ax-vasquez/wedding-site-api/settings"
)

func main() {
	settings.Setup()
	db.Setup()
	log.Println("HELLO")
}
