package web

import "net/http"

// HTTPResponse is the response interface
type HTTPResponse interface {
	CreateResponse() error
}

// Response is a response object which wrap http.ResponseWriter
type Response struct {
	w        http.ResponseWriter
	headers  map[string]string
	cookie   *http.Cookie
	original []byte
	code     int
}

// SetCode set response code
func (resp *Response) SetCode(code int) {
	resp.code = code
}

// ResponseWriter return the http.ResponseWriter
func (resp *Response) ResponseWriter() http.ResponseWriter {
	return resp.w
}

// SetContent set response content
func (resp *Response) SetContent(content []byte) {
	resp.original = content
}

// Header set response header
func (resp *Response) Header(key, value string) {
	resp.headers[key] = value
}

// Cookie set cookie
func (resp *Response) Cookie(cookie *http.Cookie) {
	// http.SetCookie(resp.w, cookie)
	resp.cookie = cookie
}

// Flush send all response contents to client
func (resp *Response) Flush() {
	// set response headers
	for key, value := range resp.headers {
		resp.w.Header().Set(key, value)
	}

	// set cookies
	if resp.cookie != nil {
		http.SetCookie(resp.w, resp.cookie)
	}

	// set response code
	resp.w.WriteHeader(resp.code)

	// send response body
	_, _ = resp.w.Write([]byte(resp.original))
}

// SendError send a error response
func (resp *Response) SendError(code int, message string) {
	resp.w.WriteHeader(code)
	_, _ = resp.w.Write([]byte(message))
}

// M represents a kv response items
type M map[string]interface{}
