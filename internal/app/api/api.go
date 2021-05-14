package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/serbi/nasa_apod_collector/internal/pkg/logger"
	"github.com/serbi/nasa_apod_collector/internal/pkg/models"
	"github.com/serbi/nasa_apod_collector/internal/pkg/settings"
)

const ComponentApi = "Component API"

func WriteJSONResponse(w http.ResponseWriter, res *models.ApiResponse) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(res.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.PrintLogError(&models.LogMessage{
			Reference: ComponentApi,
			Message:   "json.NewEncoder.Encode: " + err.Error(),
		})
		return
	}
	w.Header().Set(settings.ContentTypeHeader, settings.JsonMime)
	w.WriteHeader(res.Status)
	_, err := io.Copy(w, &buf)
	if err != nil {
		logger.PrintLogError(&models.LogMessage{
			Reference: ComponentApi,
			Message:   "io.Copy: " + err.Error(),
		})
		return
	}
	return
}
