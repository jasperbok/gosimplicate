package gosimplicate

import "fmt"

// SimplicateAPIError represents an error generated by Simplicate's API.
type SimplicateAPIError struct {
	Code              int    `json:"code"`
	Type              string `json:"type"`
	Message           string `json:"message"`
	TranslatedMessage string `json:"translated_message"`
	ErrorMessage      string `json:"error_message"`
}

func (e SimplicateAPIError) Error() string {
	return fmt.Sprintf("%d (%s): %s", e.Code, e.Type, e.Message)
}