package models

//easyjson:json
type Error struct {
	Message string `json:"message"`
}

//easyjson:json
type Errors []Error

func NewError(message string) error {
	return &Error{
		Message: message,
	}
}

func (err Error) Error() string {
	return err.Message
}

var (
	ErrDoesNotExist = NewError("entity does not exist")
	ErrAlreadyExist = NewError("entity already exist")
	ErrInvalid      = NewError("entity is invalid")
	ErrConflict     = NewError("entity conflicts with other entity")
	ErrAccessDenied = NewError("access to entity denied")
	ErrValidation   = NewError("entity validation failed")
	ErrBadForeign   = NewError("entity have bad foreign relation")
)
