package helper

const (
	ErrPlatformNotSupported = "your platform is unsupported! i can't clear terminal screen :("

	ErrDbUrlNotExist = "database URL not found"
	ErrEnvNotFound   = ".env file not found"

	ErrContactNotFound = "contact not found"

	ErrContactNameNotValid   = "name yang dimasukkan tidak valid"
	ErrContactNoTelpNotValid = "no_telp yang dimasukkan tidak valid"
	ErrContactIdNotValid     = "contact id yang dimasukkan tidak valid"
)

type AppError struct {
	Message string
}

func NewAppError(message string) *AppError {
	return &AppError{
		Message: message,
	}
}

func (e *AppError) Error() string {
	return e.Message
}
