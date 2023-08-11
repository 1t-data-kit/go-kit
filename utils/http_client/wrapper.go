package http_client

import "net/http"

type ClientWrapper func(client *http.Client, trace Trace)
type RequestWrapper func(request *http.Request, trace Trace) error
type ResponseWrapper func(response *http.Response, trace Trace) error
