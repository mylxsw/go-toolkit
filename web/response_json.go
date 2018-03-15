package web

import (
	"fmt"
	"encoding/json"
)

// JSONResponse json响应
type JSONResponse struct {
	response *Response
	original interface{}
}

// NewJSONResponse 创建JSONResponse对象
func NewJSONResponse(response *Response, res interface{}) JSONResponse {
	return JSONResponse{
		response: response,
		original: res,
	}
}

// CreateResponse 创建响应内容
func (resp JSONResponse) CreateResponse() error {
	res, err := json.Marshal(resp.original)
	if err != nil {
		err = fmt.Errorf("Json encode failed: %v [%v]", err, resp.original)

		return err
	}

	resp.response.Header("Content-Type", "application/json; charset=utf-8")
	resp.response.SetContent(string(res))

	resp.response.Flush()
	return nil
}

