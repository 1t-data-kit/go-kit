package http_client

import (
	"context"
	"testing"
)

func TestGet(t *testing.T) {
	tr := NewDefaultTrace()
	var response string
	if err := GET(
		context.TODO(),
		"https://www.baidu.com?wd=测试",
		WithTrace(tr),
		WithContentType("text/html"),
		WithResponseData(&response),
	); err != nil {
		t.Fatal(err)
	}
	t.Log(tr.Url, tr.Header, tr.Request)
	//t.Log(response)
}

func TestPOST(t *testing.T) {
	tr := NewDefaultTrace()
	var response string
	type a struct {
		Wd string `json:"wd"`
	}
	if err := POSTJson(
		context.TODO(),
		"https://www.baidu.com",
		map[string]interface{}{
			"wd": "xxx<aaa>&bbb",
			"a":  243,
		},
		//[]byte{1, 2, 3, 4, 5},
		WithTrace(tr),
		WithResponseData(&response),
	); err != nil {
		t.Fatal(err)
	}
	t.Log(tr.Url, tr.Header, tr.Request)
	t.Log(response)
}
