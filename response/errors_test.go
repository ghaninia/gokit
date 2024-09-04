package response

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errStub = errors.New("stub")

func TestServiceError_Error(t *testing.T) {
	resErr := NewServiceError(errStub, map[string]interface{}{
		"test": "test",
	}).SetType("test")

	if resErr.Error() != errStub.Error() {
		t.Errorf("Error() = %v, want %v", resErr.Error(), errStub.Error())
	}

	if resErr.GetType() != "test" {
		t.Errorf("GetType() = %v, want %v", resErr.GetType(), "test")
	}

	if resErr.GetAttributes()["test"] != "test" {
		t.Errorf("GetAttributes() = %v, want %v", resErr.GetAttributes()["test"], "test")
	}

	if resErr.Error() != errStub.Error() {
		t.Errorf("Error() = %v, want %v", resErr.Error(), errStub.Error())
	}
}

func TestServiceError_Validator(t *testing.T) {
	resErr := NewServiceError(errStub, map[string]interface{}{
		"test": "test",
	}).SetType("test")

	statusCode, _ := NewResponse(nil).WithError(resErr).EchoPure()
	assert.Equal(t, http.StatusInternalServerError, statusCode)
}

func TestServiceError_WithNativeError(t *testing.T) {
	statusCode, resp := NewResponse(nil).WithError(errStub).EchoPure()
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, errStub.Error(), resp["errors"].([]ErrorResponse)[0].Detail)
	assert.Equal(t, errStub.Error(), resp["errors"].([]ErrorResponse)[0].TypeInfo)
	assert.Equal(t, http.StatusInternalServerError, resp["errors"].([]ErrorResponse)[0].Status)
}

func TestServiceError_WithResponseError(t *testing.T) {
	resErr := NewServiceError(errStub, map[string]interface{}{
		"test": "test",
	}).SetType("test")

	statusCode, resp := NewResponse(nil).WithError(resErr).EchoPure()
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, "stub", resp["errors"].([]ErrorResponse)[0].Detail)
	assert.Equal(t, "test", resp["errors"].([]ErrorResponse)[0].TypeInfo)
	assert.Equal(t, http.StatusInternalServerError, resp["errors"].([]ErrorResponse)[0].Status)
	assert.Equal(t, "test", resp["errors"].([]ErrorResponse)[0].Attributes["test"])
}
