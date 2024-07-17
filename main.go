package main

import (
	"log"

	"github.com/ax-vasquez/wedding-site-api/models"
	"github.com/ax-vasquez/wedding-site-api/settings"
)

func main() {
	settings.Setup()
	models.Setup()
	log.Println("HELLO")
}
