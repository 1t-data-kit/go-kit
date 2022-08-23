package http_client

type trace interface {
	SetUrl(url string)
	SetHeader(header string)
	SetRequest(request string)
	SetResponse(response string)
	SetHttpCode(code int)
}
