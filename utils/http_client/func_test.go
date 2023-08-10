package http_client

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func TestGet(t *testing.T) {
	tr := NewDefaultTrace()
	var response string
	if err := GET(
		context.TODO(),
		"http://mashineadmin.unionscreens.com/api/device/v1/queryDevices3?pageNum=2&pageSize=2",
		WithTrace(tr),
		WithContentType("text/html"),
		WithResponseData(&response),
	); err != nil {
		t.Fatal(err)
	}
	t.Log(tr.Url, tr.Header, tr.Request)
	t.Log(response)
}

func TestPOST(t *testing.T) {
	tr := NewDefaultTrace()
	var response string
	type a struct {
		Wd string `json:"wd"`
	}
	if err := POSTJson(
		context.TODO(),
		"https://www.baidu.com?a=2",
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
