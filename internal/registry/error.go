package registry

const (
	ApplicationNotFoundCode = -1
	InstanceNotFoundCode    = -2
)

var (
	ApplicationNotFoundError = NewRegistryError(ApplicationNotFoundCode, "application not found error")
	InstanceNotFoundError    = NewRegistryError(InstanceNotFoundCode, "instance not found error")
)

type RegistryError struct {
	Code    int32
	Message string
}

func NewRegistryError(code int32, message string) RegistryError {

	return RegistryError{
		Code:    code,
		Message: message,
	}
}

func (e RegistryError) Error() string {
	return e.Message
}
