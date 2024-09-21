package errors

const (
	// Server The operation looked fine but the server couldn't process it.
	Server string = "server"
	// NotFound The requested resource could not be found.
	NotFound string = "not_found"
	// User The request was well-formed but was unable to be followed due to semantic errors.
	User string = "user"
)

// Error Representation of errors in API.
type Error struct {
	Type    string
	Message string `json:"message"`
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// UnauthorizedError Used to represent unauthorized errors.
func UnauthorizedError(message string) *Error {
	return &Error{
		Type:    User,
		Message: message,
	}
}

// MissingTenantError Used to represent missing tenant errors.
func MissingTenantError(message string) *Error {
	return &Error{
		Type:    User,
		Message: message,
	}
}

// ValidationError Used to represent validation errors.
func ValidationError(message string, err error) *Error {
	return &Error{
		Type:    User,
		Message: message,
		Err:     err,
	}
}

// NotFoundError Used to represent not found errors.
func NotFoundError(message string) *Error {
	return &Error{
		Type:    NotFound,
		Message: message,
	}
}

// ServerError Used to represent server errors.
func ServerError(message string) *Error {
	return &Error{
		Type:    Server,
		Message: message,
	}
}

// EntitlementError Used to represent entitlement errors.
func EntitlementError(message string, err error) *Error {
	return &Error{
		Type:    User,
		Message: message,
		Err:     err,
	}
}
