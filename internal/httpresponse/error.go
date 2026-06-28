package httpresponse

type ErrorResponse struct {
	Success   bool              `json:"success"`
	Code      int               `json:"code"`
	Message   string            `json:"message"`
	Errors    map[string]string `json:"errors,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

func NewErrorResponse(response ErrorResponse) *ErrorResponse {
	return &response
}