package response

import (
	"errors"
	"ghaninia/gokit/translation"

	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	errTypeMsgSomethingIsWrong  = "server.errors.something_is_wrong"
	errTypeInfoSomethingIsWrong = "something_is_wrong"
)

type Response interface {
	Validation(err error) Response
	WithPayload(data any) Response
	WithError(err error) Response
	WithMeta(data interface{}) Response
	Echo(ctx *gin.Context)
	EchoPure() (statusCode int, response map[string]any)
}

type response struct {
	statusCodeMapping map[string]int
	translation       translation.Translation
	response          map[string]interface{}
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

	return &response{
		statusCodeMapping: statusCodeMapping,
		translation:       trans,
		response:          make(map[string]interface{}),
	}
}

// Validation sets the validation error to be sent to the client.
func (r *response) Validation(err error) Response {
	v := newValidationTranslator(r.translation).translate(err)
	r.validation = &v
	return r
}

// WithMeta sets the meta data to be sent to the client.
func (r *response) WithMeta(data interface{}) Response {
	r.response["meta"] = data
	return r
}

// WithPayload sets the data to be sent to the client.
func (r *response) WithPayload(data any) Response {
	r.payload = &data
	return r
}

// WithError sets the error to be sent to the client.
func (r *response) WithError(err error) Response {
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
func (r *response) WithStatusCode(statusCode int) Response {
	r.statusCode = &statusCode
	return r
}

// EchoPure returns the response to be sent to the client.
func (r *response) EchoPure() (statusCode int, response map[string]any) {
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

	return statusCode, r.response
}

// Echo sends the response to the client.
func (r *response) Echo(ctx *gin.Context) {
	statusCode, rsp := r.EchoPure()
	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		ctx.JSON(statusCode, rsp)
		return
	}
	ctx.AbortWithStatusJSON(statusCode, rsp)
}

// getStatusMapping returns the status code based on the error message.
// If the error message is not found in the status code mapping, it returns 500.
func (r *response) getStatusMapping() (statusCode int) {
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
