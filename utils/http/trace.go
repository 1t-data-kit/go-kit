package http

import "encoding/json"

type Trace struct {
	Url      string `json:"url"`
	Header   interface{} `json:"header"`
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
	HttpCode int    `json:"httpCode"`
}

func NewTrace() *Trace {
	return &Trace{}
}

func (trace *Trace) String() string {
	data,_ := json.Marshal(trace)
	return string(data)
}
