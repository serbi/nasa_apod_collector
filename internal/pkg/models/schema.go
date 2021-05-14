package models

const (
	ErrorMsgPublic500 ApiErrorMessage = "An internal error occurred while processing your request"
	ErrorMsgPublic422 ApiErrorMessage = "The request parameters were unprocessable"
	ErrorMsgPublic404 ApiErrorMessage = "The requested resource does not exist"
	ErrorMsgPublic405 ApiErrorMessage = "This HTTP Method is not allowed"
)
