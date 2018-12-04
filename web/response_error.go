package web

// ErrorResponse Error Response
type ErrorResponse struct {
	response *Response
	original string
	code     int
}

// NewErrorResponse Create error response
func NewErrorResponse(response *Response, res string, code int) ErrorResponse {
	return ErrorResponse{
		response: response,
		original: res,
		code:     code,
	}
}

// CreateResponse 创建响应内容
func (resp ErrorResponse) CreateResponse() error {
	resp.response.SendError(resp.code, resp.original)
	resp.response.Flush()
	return nil
}
