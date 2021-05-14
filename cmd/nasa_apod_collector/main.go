package main

import (
	"log"
	"net/http"

	"github.com/serbi/nasa_apod_collector/internal/app/router"
	"github.com/serbi/nasa_apod_collector/internal/pkg/settings"
)

func main() {
	router.NewRouter().RegisterHandlers()

	port := settings.AppPort
	log.Printf("Listening on %s ...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
