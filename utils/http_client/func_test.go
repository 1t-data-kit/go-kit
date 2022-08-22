package http_client

import (
	"context"
	"testing"
)

func TestGet(t *testing.T) {
	trace := NewTrace()
	var response string
	if err := GET(
		context.TODO(),
		"https://www.baidu.com?wd=测试",
		WithTrace(trace),
		WithContentType("text/html"),
		WithResponseData(&response),
	); err != nil {
		t.Fatal(err)
	}
	t.Log(trace.Url, trace.Header, trace.Request)
	//t.Log(response)
}

func TestPOST(t *testing.T) {
	trace := NewTrace()
	var response string
	type a struct {
		Wd string `json:"wd"`
	}
	if err := POSTBinary(
		context.TODO(),
		"https://www.baidu.com",
		//map[string]interface{}{
		//	"wd": "xxx",
		//	"a":  243,
		//},
		//map[string][]byte{
		//	"file1": []byte{1, 2, 3, 4},
		//	"file2": []byte{3, 4, 5},
		//},
		[]byte{1, 2, 3, 4, 5},
		WithTrace(trace),
		WithResponseData(&response),
	); err != nil {
		t.Fatal(err)
	}
	t.Log(trace.Header, trace.Request)
	t.Log(response)
}
