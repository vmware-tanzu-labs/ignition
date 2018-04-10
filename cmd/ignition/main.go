package main

import (
	"log"

	"github.com/pivotalservices/ignition/config"
	"github.com/pivotalservices/ignition/http"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.LUTC)
	ignition, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	api := http.API{
		Ignition: ignition,
	}
	log.Printf("Starting Server listening on %s\n", api.URI())
	log.Fatal(api.Run())
}
