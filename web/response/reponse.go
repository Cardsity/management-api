package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Every response from the api contains this.
type BasicApiResponse struct {
	ErrorOccurred bool `json:"err"`
}

// Every error response from the api is in this format.
type BasicApiResponseError struct {
	BasicApiResponse
	Errors []string `json:"errors"`
}

// Returns an instance of BasicApiResponseError which is with the supplied errors.
func NewApiErrorResponse(errors []string) BasicApiResponseError {
	return BasicApiResponseError{
		BasicApiResponse: BasicApiResponse{
			ErrorOccurred: true,
		},
		Errors: errors,
	}
}

// Every response that is ok/was successful from the api is in this format.
type BasicApiResponseOk struct {
	BasicApiResponse
	Data interface{} `json:"data"`
}

// Returns an instance of BasicApiResponseOk which is with the supplied parameters.
func NewApiOkResponse(data interface{}) BasicApiResponseOk {
	return BasicApiResponseOk{
		BasicApiResponse: BasicApiResponse{
			ErrorOccurred: false,
		},
		Data: data,
	}
}

func BadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, NewApiErrorResponse([]string{ErrorBadRequest}))
}

func Conflict(c *gin.Context, reason string) {
	c.JSON(http.StatusConflict, NewApiErrorResponse([]string{reason}))
}

func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, NewApiErrorResponse([]string{ErrorInternal}))
}

func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, NewApiOkResponse(data))
}
