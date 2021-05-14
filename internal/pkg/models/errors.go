package models

type ApiErrorMessage string

func (s ApiErrorMessage) String() string {
	return string(s)
}

type ApiError struct {
	Message ApiErrorMessage `json:"error"`
}

func NewApiError(message ApiErrorMessage) (apiErr *ApiError) {
	apiErr = &ApiError{Message: message}
	return
}

func (a ApiError) Error() string {
	return a.Message.String()
}
