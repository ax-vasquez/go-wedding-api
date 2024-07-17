package main

import (
	"log"

	"github.com/ax-vasquez/wedding-site-api/settings"
)

func main() {
	settings.Setup()
	log.Println("HELLO")
}
