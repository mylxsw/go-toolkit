package web


// HTMLResponse HTML响应
type HTMLResponse struct {
	response *Response
	original string
}

// NewHTMLResponse 创建HTML响应
func NewHTMLResponse(response *Response, res string) HTMLResponse {
	return HTMLResponse{
		response: response,
		original: res,
	}
}

// CreateResponse 创建响应内容
func (resp HTMLResponse) CreateResponse() error {
	resp.response.Header("Content-Type", "text/html; charset=utf-8")
	resp.response.SetContent(resp.original)

	resp.response.Flush()
	return nil
}

