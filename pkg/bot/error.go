package botlogic

import (
	"encoding/json"
	"fmt"
)

type CustomError struct {
	ChatID int64
	URL    string
	Params map[string]interface{}
	Err    error
}

// Implement the error interface
func (e *CustomError) Error() string {
	paramsJSON, _ := json.Marshal(e.Params) // Convert params to JSON string
	return fmt.Sprintf("ChatId: %d | URL: %s | Params: %s | Cause: %v",
		e.ChatID, e.URL, string(paramsJSON), e.Err)
}

// Unwrap allows errors.Is and errors.As to work
func (e *CustomError) Unwrap() error {
	return e.Err
}
