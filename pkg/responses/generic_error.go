package responses

type GenericError struct {
	Error string `json:"error"`
}

func NewGenericError(msg string) *GenericError {
	return &GenericError{Error: msg}
}
