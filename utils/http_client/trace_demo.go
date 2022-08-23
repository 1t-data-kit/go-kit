package http_client

import (
	"encoding/json"
)

type TraceDemo struct {
	Url      string `json:"url"`
	Header   string `json:"header"`
	Request  string `json:"request"`
	Response string `json:"response"`
	HttpCode int    `json:"httpCode"`
}

func NewTraceDemo() *TraceDemo {
	return &TraceDemo{}
}

func (demo *TraceDemo) String() string {
	data, _ := json.Marshal(demo)
	return string(data)
}

func (demo *TraceDemo) SetUrl(url string) {
	demo.Url = url
}

func (demo *TraceDemo) SetHeader(header string) {
	demo.Header = header
}

func (demo *TraceDemo) SetRequest(request string) {
	demo.Request = request
}

func (demo *TraceDemo) SetResponse(response string) {
	demo.Response = response
}

func (demo *TraceDemo) SetHttpCode(code int) {
	demo.HttpCode = code
}
