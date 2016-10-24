package goe


type HttpError struct {
	Code int
	Msg  string
}

func (e HttpError) Error() string {
	return e.Msg
}

func NewHttpError(code int, msg string) *HttpError {
	return &HttpError{
		Code:code,
		Msg:msg,
	}
}