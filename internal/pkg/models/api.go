package models

type ApiResponse struct {
	Status int
	Body   *ApiResponseBody
}

type ApiResponseBody struct {
	*ApiUrls
	*ApiError
}

type ApiUrls struct {
	Urls []string `json:"urls"`
}
