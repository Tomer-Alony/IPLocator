package errors

type HttpError struct {
	Message string `json:"error"`
	Code int `json:"code"`
}

func (h HttpError) Error() string {
	return h.Message
}
