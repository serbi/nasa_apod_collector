package router

import (
	"net/http"

	"github.com/gorilla/schema"

	"github.com/serbi/nasa_apod_collector/internal/app/api"
	"github.com/serbi/nasa_apod_collector/internal/app/controllers"
	"github.com/serbi/nasa_apod_collector/internal/pkg/logger"
	"github.com/serbi/nasa_apod_collector/internal/pkg/models"
)

const ComponentRouterHandlers = "Component Router Handlers"

func HandlePictures(w http.ResponseWriter, r *http.Request) {
	res := &models.ApiResponse{}
	if r.Method != "GET" {
		res = &models.ApiResponse{
			Status: http.StatusMethodNotAllowed,
			Body: &models.ApiResponseBody{
				ApiError: models.NewApiError(models.ErrorMsgPublic405),
			},
		}
		logger.PrintLogError(&models.LogMessage{
			Reference: ComponentRouterHandlers,
			Message:   "r.Method != 'GET'",
		})
		api.WriteJSONResponse(w, res)
		return
	}

	decoder := schema.NewDecoder()
	var params models.PicturesGetParams
	err := decoder.Decode(&params, r.URL.Query())
	if err != nil {
		res = &models.ApiResponse{
			Status: http.StatusInternalServerError,
			Body: &models.ApiResponseBody{
				ApiError: models.NewApiError(models.ErrorMsgPublic500),
			},
		}
		logger.PrintLogError(&models.LogMessage{
			Reference: ComponentRouterHandlers,
			Message:   "decoder.Decode: " + err.Error(),
		})
		api.WriteJSONResponse(w, res)
		return
	}
	if len(params.StartDate) < 1 || len(params.EndDate) < 1 {
		res = &models.ApiResponse{
			Status: http.StatusUnprocessableEntity,
			Body: &models.ApiResponseBody{
				ApiError: models.NewApiError(models.ErrorMsgPublic422),
			},
		}
		logger.PrintLogError(&models.LogMessage{
			Reference: ComponentRouterHandlers,
			Message:   "len(params.StartDate) < 1 || len(params.EndDate) < 1",
		})
		api.WriteJSONResponse(w, res)
		return
	}

	res = controllers.GetPictures(&params)

	api.WriteJSONResponse(w, res)
	return
}
