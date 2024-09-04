package response

type Error interface {
	Error() string
	GetType() string
	GetAttributes() map[string]interface{}
}

type ServiceError struct {
	err        error
	errorType  string
	attributes map[string]interface{}
}

func NewServiceError(
	err error,
	attrs ...map[string]interface{},
) *ServiceError {
	attributes := make(map[string]interface{})

	if len(attrs) > 0 {
		attributes = attrs[0]
	}

	return &ServiceError{
		err:        err,
		attributes: attributes,
	}
}

func (e *ServiceError) SetType(errorType string) Error {
	e.errorType = errorType
	return e
}

func (e *ServiceError) GetMessage() string {
	return e.err.Error()
}

func (e *ServiceError) GetType() string {
	return e.errorType
}

func (e *ServiceError) GetAttributes() map[string]interface{} {
	return e.attributes
}

func (e *ServiceError) Error() string {
	return e.err.Error()
}
