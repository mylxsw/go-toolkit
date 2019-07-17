package next

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mylxsw/asteria/log"
)

var logger = log.Module("next")

// CreateHTTPHandler create a http handler for request processing
func CreateHTTPHandler(config *Config) http.Handler {
	rootDir := filepath.Dir(config.EndpointFile)

	handler := Handler{
		Rules:           config.Rules,
		Root:            rootDir,
		FileSys:         http.Dir(rootDir),
		SoftwareName:    config.SoftwareName,
		SoftwareVersion: config.SoftwareVersion,
		ServerName:      config.ServerIP,
		ServerPort:      strconv.Itoa(config.ServerPort),
	}

	return &HTTPHandler{
		handler: handler,
		config:  config,
	}
}

// RequestLogHandler request log handler func
type RequestLogHandler func(rc *RequestContext)
type ErrorResponseHandler func(w http.ResponseWriter, r *http.Request, code int, err error)

// Config config object for create a handler
type Config struct {
	EndpointFile         string
	ServerIP             string
	ServerPort           int
	SoftwareName         string
	SoftwareVersion      string
	RequestLogHandler    RequestLogHandler
	ErrorResponseHandler ErrorResponseHandler
	Rules                []Rule
}

// HTTPHandler http request handler wrapper
type HTTPHandler struct {
	handler Handler
	config  *Config
}

// RequestContext request context information
type RequestContext struct {
	UA      string
	Method  string
	Referer string
	Headers http.Header
	URI     string
	Body    url.Values
	Consume float64
	Code    int
	Error   string
}

// ToMap convert the requestContext to a map
func (rc *RequestContext) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"ua":      rc.UA,
		"method":  rc.Method,
		"referer": rc.Referer,
		"headers": rc.Headers,
		"uri":     rc.URI,
		"body":    rc.Body,
		"consume": rc.Consume,
		"code":    rc.Code,
		"error":   rc.Error,
	}
}

// ServeHTTP implements http.Handler interface
func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var statusCode = 200
	respWriter := NewResponseWriter(w, func(code int) {
		statusCode = code
	})

	var err error
	defer func(startTime time.Time) {
		consume := time.Now().Sub(startTime)
		if r.Form == nil {
			r.ParseForm()
		}

		if h.config.RequestLogHandler != nil {
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logger.Errorf("request log handler has some error: %v", err)
					}
				}()

				errorMsg := ""
				if err != nil {
					errorMsg = err.Error()
				}

				h.config.RequestLogHandler(&RequestContext{
					UA:      r.UserAgent(),
					Method:  r.Method,
					Referer: r.Referer(),
					Headers: r.Header,
					URI:     r.RequestURI,
					Body:    r.Form,
					Consume: consume.Seconds(),
					Code:    statusCode,
					Error:   errorMsg,
				})
			}()
		}

	}(time.Now())

	code, err := h.handler.ServeHTTP(respWriter, r)
	if err != nil {
		logger.Errorf("request failed, code=%d, err=%s", code, err.Error())
	}

	if code != 0 {
		respWriter.WriteHeader(code)
		if h.config.ErrorResponseHandler != nil {
			h.config.ErrorResponseHandler(respWriter, r, code, err)
		}
	}
}
