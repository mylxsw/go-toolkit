package next

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mylxsw/asteria/log"
)

// CreateHTTPHandler create a http handler for request processing
func CreateHTTPHandler(config *Config) http.Handler {
	rootDir := filepath.Dir(config.EndpointFile)

	handler := Handler{
		Rules:              config.Rules,
		Root:               rootDir,
		FileSys:            http.Dir(rootDir),
		SoftwareName:       config.SoftwareName,
		SoftwareVersion:    config.SoftwareVersion,
		ServerName:         config.ServerIP,
		ServerPort:         strconv.Itoa(config.ServerPort),
		OverrideHostHeader: config.OverrideHostHeader,
	}

	return &HTTPHandler{
		handler:         handler,
		config:          config,
		errorLogHandler: config.ErrorLogHandler,
	}
}

// RequestLogHandler request log handler func
type RequestLogHandler func(rc *RequestContext)
type ErrorResponseHandler func(w http.ResponseWriter, r *http.Request, code int, err error)
type ErrorLogHandler func(err error)

// Config config object for create a handler
type Config struct {
	EndpointFile         string
	ServerIP             string
	ServerPort           int
	SoftwareName         string
	SoftwareVersion      string
	RequestLogHandler    RequestLogHandler
	ErrorResponseHandler ErrorResponseHandler
	ErrorLogHandler      ErrorLogHandler
	Rules                []Rule
	OverrideHostHeader   string
}

// HTTPHandler http request handler wrapper
type HTTPHandler struct {
	handler         Handler
	errorLogHandler ErrorLogHandler
	config          *Config
}

// RequestContext request context information
type RequestContext struct {
	UA      string
	Method  string
	Referer string
	Headers http.Header
	URI     string
	Body    string
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

	body, _ := ioutil.ReadAll(r.Body)
	_ = r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var err error
	defer func(startTime time.Time) {
		consume := time.Now().Sub(startTime)
		if h.config.RequestLogHandler != nil {
			go func() {
				defer func() {
					if err := recover(); err != nil {
						log.Errorf("request log handler has some error: %v", err)
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
					Body:    string(body),
					Consume: consume.Seconds(),
					Code:    statusCode,
					Error:   errorMsg,
				})
			}()
		}

	}(time.Now())

	code, err := h.handler.ServeHTTP(respWriter, r)
	if err != nil {
		if code == 0 && h.errorLogHandler != nil {
			h.errorLogHandler(err)
		} else {
			log.Errorf("request failed, code=%d, err=%s", code, err.Error())
		}
	}

	if code != 0 {
		respWriter.WriteHeader(code)
		if h.config.ErrorResponseHandler != nil {
			h.config.ErrorResponseHandler(respWriter, r, code, err)
		}
	}
}
