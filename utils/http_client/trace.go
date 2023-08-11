package http_client

import "reflect"

type Trace interface {
	SetUrl(url string)
	SetHeader(header string)
	SetRequest(request string)
	SetResponse(response string)
	SetHttpCode(code int)
}

func traceIsNil(trace Trace) bool {
	if trace == nil || reflect.ValueOf(trace).IsNil() {
		return true
	}
	return false
}
