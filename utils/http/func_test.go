package http

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
		WithResponseData(&response),
	); err != nil {
		t.Fatal(err)
	}
	t.Log(trace)
	t.Log(response)
}

func TestPOSTJson(t *testing.T) {
	trace := NewTrace()
	var response string
	type a struct {
		Wd string `json:"wd"`
	}
	if err := POSTJson(
		context.TODO(),
		"https://www.baidu.com",
		&a{
			Wd: "xxx",
		},
		WithTrace(trace),
		WithResponseData(&response),
	); err != nil {
		t.Fatal(err)
	}
	t.Log(trace.HttpCode)
	t.Log(response)
}
