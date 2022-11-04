package http_client

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func GET(ctx context.Context, url string, option ...Option) error {
	return do(ctx, "GET", url, option...)
}

func POST(ctx context.Context, url string, source map[string]interface{}, option ...Option) error {
	return do(ctx, "POST", url, append(option, WithBody(source))...)
}

func POSTJson(ctx context.Context, url string, source interface{}, option ...Option) error {
	return do(ctx, "POST", url, append(option, WithJSONBody(source))...)
}

func POSTJsonWithEscape(ctx context.Context, url string, source interface{}, option ...Option) error {
	return do(ctx, "POST", url, append(option, WithEscapeJSONBody(source))...)
}

func POSTBinary(ctx context.Context, url string, file []byte, option ...Option) error {
	return do(ctx, "POST", url, append(option, WithBinaryBody(file))...)
}

func POSTMultipart(ctx context.Context, url string, source map[string]interface{}, file map[string][]byte, option ...Option) error {
	return do(ctx, "POST", url, append(option, WithMultipartBody(source, file))...)
}

func getString(data interface{}) string {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Struct, reflect.Slice, reflect.Map:
		m, _ := json.Marshal(data)
		return string(m)
	}
	return ""
}

func getStrings(sep string, data ...interface{}) string {
	stringList := make([]string, 0)
	for _, d := range data {
		stringList = append(stringList, getString(d))
	}
	return strings.Join(stringList, sep)
}

func do(ctx context.Context, method string, url string, option ...Option) error {
	_optionList := optionList(option)

	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return err
	}
	if err = _optionList.Request(request); err != nil {
		return err
	}

	response, err := _optionList.Client().Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err = _optionList.Response(response); err != nil {
		return err
	}

	return nil
}
