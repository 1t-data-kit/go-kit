package http_client

import (
	"encoding/json"
)

type DefaultTrace struct {
	Url      string `json:"url"`
	Header   string `json:"header"`
	Request  string `json:"request"`
	Response string `json:"response"`
	HttpCode int    `json:"httpCode"`
}

func NewDefaultTrace() *DefaultTrace {
	return &DefaultTrace{}
}

func (defaultTrace *DefaultTrace) String() string {
	data, _ := json.Marshal(defaultTrace)
	return string(data)
}

func (defaultTrace *DefaultTrace) SetUrl(url string) {
	defaultTrace.Url = url
}

func (defaultTrace *DefaultTrace) SetHeader(header string) {
	defaultTrace.Header = header
}

func (defaultTrace *DefaultTrace) SetRequest(request string) {
	defaultTrace.Request = request
}

func (defaultTrace *DefaultTrace) SetResponse(response string) {
	defaultTrace.Response = response
}

func (defaultTrace *DefaultTrace) SetHttpCode(code int) {
	defaultTrace.HttpCode = code
}
