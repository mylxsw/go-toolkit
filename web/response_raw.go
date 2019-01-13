package web

// RawResponse 原生响应
type RawResponse struct {
	response *Response
}

// NewRawResponse create a RawResponse
func NewRawResponse(response *Response) *RawResponse {
	return &RawResponse{response: response}
}

// Response get real response object
func (resp *RawResponse) Response() *Response {
	return resp.response
}

// CreateResponse flush response to client
func (resp *RawResponse) CreateResponse() error {
	resp.response.Flush()
	return nil
}
