package httpresponse

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

func NewSuccessResponse(response SuccessResponse) *SuccessResponse {
	return &response
}
