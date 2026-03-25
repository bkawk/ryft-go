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
	var parsed APIError
	parsed.Status = status

	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err == nil {
		if value, ok := raw["code"].(string); ok {
			parsed.Code = value
		}
		if value, ok := raw["message"].(string); ok {
			parsed.Message = value
		}
		if value, ok := raw["requestId"].(string); ok {
			parsed.RequestID = value
		}
		if value, ok := raw["errors"]; ok {
			encoded, marshalErr := json.Marshal(value)
			if marshalErr == nil {
				_ = json.Unmarshal(encoded, &parsed.Errors)
			}
		}
	}

	if parsed.Message == "" {
		parsed.Message = string(body)
	}

	return &parsed
}
