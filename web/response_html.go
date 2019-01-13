package web

// HTMLResponse HTML响应
type HTMLResponse struct {
	response *Response
	original string
	code     int
}

// NewHTMLResponse 创建HTML响应
func NewHTMLResponse(response *Response, code int, res string) *HTMLResponse {
	return &HTMLResponse{
		response: response,
		original: res,
		code:     code,
	}
}

// WithCode set response code and return itself
func (resp *HTMLResponse) WithCode(code int) *HTMLResponse {
	resp.code = code
	return resp
}

// CreateResponse 创建响应内容
func (resp *HTMLResponse) CreateResponse() error {
	resp.response.SetCode(resp.code)
	resp.response.Header("Content-Type", "text/html; charset=utf-8")
	resp.response.SetContent([]byte(resp.original))

	resp.response.Flush()
	return nil
}
