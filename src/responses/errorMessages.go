package responses

const (
	ErrMsgBadRequest          = "the request could not be understood or was missing required parameters"
	ErrMsgUnauthorized        = "authentication is required and has failed or has not yet been provided"
	ErrMsgForbidden           = "you do not have permission to access this resource"
	ErrMsgNotFound            = "the requested resource was not found"
	ErrMsgUnprocessableEntity = "the request could not be processed"
	ErrMsgInternalServerError = "an unexpected error occurred. Please try again later"
)
