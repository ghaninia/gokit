package response

import (
	"errors"
	"github.com/ghaninia/gokit/meta"
	"github.com/ghaninia/gokit/translation"

	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	errTypeMsgSomethingIsWrong  = "server.errors.something_is_wrong"
	errTypeInfoSomethingIsWrong = "something_is_wrong"
)

type NormalizeResponse struct {
	Data    *interface{} `json:"data;omitempty"`
	Message *string      `json:"message;omitempty"`
	Errors  *interface{} `json:"errors;omitempty"`
	Meta    *meta.Meta   `json:"meta;omitempty"`
}

type Response interface {
	Validation(err error) Response
	WithPayload(data any) Response
	WithMessage(message string, args ...map[string]interface{}) Response
	WithError(err error) Response
	WithMeta(data interface{}) Response
	Echo(ctx *gin.Context)
	EchoPure() (statusCode int, response map[string]any)
	WithStatusCode(statusCode int) Response
}

type Resource struct {
	statusCodeMapping map[string]int
	translation       translation.Translation
	response          map[string]interface{}
	message           *string
	payload           *any
	validation        *Validations
	statusCode        *int
	nativeError       error
	responseError     Error
}

type ErrorResponse struct {
	TypeInfo   string                 `json:"type_info"`
	Status     int                    `json:"status"`
	Detail     string                 `json:"detail"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// NewResponse creates a new response.
func NewResponse(
	trans translation.Translation,
	statusCodeMappings ...map[string]int,
) Response {
	statusCodeMapping := make(map[string]int)

	if len(statusCodeMappings) > 0 {
		statusCodeMapping = statusCodeMappings[0]
	}

	return &Resource{
		statusCodeMapping: statusCodeMapping,
		translation:       trans,
		response:          make(map[string]interface{}),
	}
}

// Validation sets the validation error to be sent to the client.
func (r *Resource) Validation(err error) Response {
	v := newValidationTranslator(r.translation).translate(err)
	r.validation = &v
	return r
}

// WithMessage sets the message to be sent to the client.
func (r *Resource) WithMessage(message string, args ...map[string]interface{}) Response {
	var arg map[string]interface{}
	if len(args) > 0 {
		arg = args[0]
	}
	trans := r.translation.Trans(message, arg)
	r.message = &trans
	return r
}

// WithMeta sets the meta data to be sent to the client.
func (r *Resource) WithMeta(data interface{}) Response {
	r.response["meta"] = data
	return r
}

// WithPayload sets the data to be sent to the client.
func (r *Resource) WithPayload(data any) Response {
	r.payload = &data
	return r
}

// WithError sets the error to be sent to the client.
func (r *Resource) WithError(err error) Response {
	if err == nil {
		return r
	}

	var e Error
	if errors.As(err, &e) {
		r.responseError = e
	} else {
		r.nativeError = err
	}

	return r
}

// WithStatusCode sets the status code to be sent to the client.
func (r *Resource) WithStatusCode(statusCode int) Response {
	r.statusCode = &statusCode
	return r
}

// EchoPure returns the response to be sent to the client.
func (r *Resource) EchoPure() (statusCode int, response map[string]any) {
	if r.statusCode != nil {
		statusCode = *r.statusCode
	} else {
		statusCode = r.getStatusMapping()
	}

	errInfo := errTypeInfoSomethingIsWrong
	errDetail := errTypeMsgSomethingIsWrong
	errAttributes := make(map[string]interface{})

	if r.nativeError != nil {
		if r.nativeError.Error() != "" {
			errInfo = r.nativeError.Error()
			errDetail = r.nativeError.Error()
		}
	} else if r.responseError != nil {
		if r.responseError.GetType() != "" {
			errInfo = r.responseError.GetType()
		}
		if r.responseError.Error() != "" {
			errDetail = r.responseError.Error()
		}
		errAttributes = r.responseError.GetAttributes()
	}

	if r.translation != nil {
		errDetail = r.translation.Trans(errDetail, errAttributes)
	}

	if r.nativeError != nil || r.responseError != nil {
		r.response["errors"] = []ErrorResponse{
			{
				TypeInfo:   errInfo,
				Status:     statusCode,
				Detail:     errDetail,
				Attributes: errAttributes,
			},
		}
	}

	if r.validation != nil {
		r.response["errors"] = *r.validation
	}

	if r.payload != nil {
		r.response["data"] = *r.payload
	}

	if r.statusCode != nil {
		statusCode = *r.statusCode
	}

	if r.message != nil {
		r.response["message"] = *r.message
	}

	return statusCode, r.response
}

// Echo sends the response to the client.
func (r *Resource) Echo(ctx *gin.Context) {
	statusCode, rsp := r.EchoPure()
	response := NormalizeResponse{
		Data: func() *interface{} {
			if rsp["data"] == nil {
				return nil
			}
			res := rsp["data"]
			return &res
		}(),
		Message: func() *string {
			if rsp["message"] == nil {
				return nil
			}
			res := rsp["message"].(string)
			return &res
		}(),
		Errors: func() *interface{} {
			if rsp["errors"] == nil {
				return nil
			}
			res := rsp["errors"]
			return &res
		}(),
		Meta: func() *meta.Meta {
			if rsp["meta"] == nil {
				return nil
			}
			res := rsp["meta"].(meta.Meta)
			return &res
		}(),
	}
	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		ctx.JSON(statusCode, response)
		return
	}
	ctx.AbortWithStatusJSON(statusCode, response)
}

// getStatusMapping returns the status code based on the error message.
// If the error message is not found in the status code mapping, it returns 500.
func (r *Resource) getStatusMapping() (statusCode int) {
	switch {
	case r.responseError != nil:
		{
			e := r.responseError
			msg := e.Error()
			if msg == "" {
				return http.StatusInternalServerError
			}

			if val, ok := r.statusCodeMapping[msg]; !ok {
				statusCode = http.StatusInternalServerError
			} else {
				statusCode = val
			}
		}
	case r.nativeError != nil:
		{
			statusCode = http.StatusInternalServerError
		}
	default:
		{
			statusCode = http.StatusOK
		}
	}

	return statusCode
}
