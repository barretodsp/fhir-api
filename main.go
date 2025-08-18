package main

import (
	"log"

	_ "fhir-api/docs"
)

func main() {
	app := RunApp()
	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
