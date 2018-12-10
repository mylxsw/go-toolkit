package web

// RawResponse 原生响应
type RawResponse struct {
	response *Response
}

// NewRawResponse create a RawResponse
func NewRawResponse(response *Response) RawResponse {
	return RawResponse{response: response}
}

func (resp RawResponse) CreateResponse() error {
	resp.response.Flush()
	return nil
}
