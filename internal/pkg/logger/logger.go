package logger

import (
	"log"

	"github.com/serbi/nasa_apod_collector/internal/pkg/models"
)

func PrintLogError(l *models.LogMessage) {
	log.Print("ERROR in " + l.Reference + ": " + l.Message)
	return
}

func PrintLog(l *models.LogMessage) {
	log.Print(l.Reference + ": " + l.Message)
	return
}
