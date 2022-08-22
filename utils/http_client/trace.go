package http_client

import (
	"encoding/json"
	"net/http"
)

type Trace struct {
	Url      string      `json:"url"`
	Header   interface{} `json:"header"`
	Request  string      `json:"request"`
	Response string      `json:"response"`
	HttpCode int         `json:"httpCode"`
}

func NewTrace() *Trace {
	return &Trace{}
}

func (trace *Trace) String() string {
	data, _ := json.Marshal(trace)
	return string(data)
}

func (trace *Trace) Record(source interface{}, data ...interface{}) {
	if request, ok := source.(*http.Request); ok {
		trace.Url = request.URL.String()
		trace.Header = request.Header
		trace.Request = getStrings(",", data)
	}
	if response, ok := source.(*http.Response); ok {
		trace.HttpCode = response.StatusCode
		trace.Response = getStrings(",", data)
	}
}
