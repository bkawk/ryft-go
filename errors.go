package ryft

import "encoding/json"

type APIError struct {
	Status    int              `json:"status"`
	Code      string           `json:"code,omitempty"`
	Message   string           `json:"message,omitempty"`
	RequestID string           `json:"requestId,omitempty"`
	Errors    []APIErrorDetail `json:"errors,omitempty"`
}

type APIErrorDetail struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Field   string `json:"field,omitempty"`
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Code != "" {
		return e.Code
	}
	return "ryft api error"
}

func parseAPIError(status int, body []byte) error {
	parsed := APIError{Status: status}
	_ = json.Unmarshal(body, &parsed)

	if parsed.Message == "" {
		parsed.Message = string(body)
	}

	return &parsed
}
