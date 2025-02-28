package httpError

import (
	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) SendError(c *gin.Context) {
	if e.StatusCode == 0 {
		e.StatusCode = 500
	}
	c.JSON(e.StatusCode, gin.H{"status": "error", "message": e.Error()})
}

func New(statusCode int, message string) *HTTPError {
	return &HTTPError{
		Message:    message,
		StatusCode: statusCode,
	}
}
