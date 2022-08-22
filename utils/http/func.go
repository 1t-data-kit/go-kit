package http

import (
	"context"
	"net/http"
	"strconv"
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

func POSTBinary(ctx context.Context, url string, file []byte, option ...Option) error {
	return do(ctx, "POST", url, append(option, WithBinaryBody(file))...)
}

func POSTMultipart(ctx context.Context, url string, file map[string][]byte, source map[string]interface{}, option ...Option) error {
	return do(ctx, "POST", url, append(option, WithMultipartBody(file, source))...)
}

func getAssertString(data interface{}) string {
	if str, ok := data.(string); ok {
		return str
	} else if intValue, ok := data.(int); ok {
		return strconv.Itoa(intValue)
	} else if intValue, ok := data.(int32); ok {
		return strconv.FormatInt(int64(intValue), 10)
	} else if intValue, ok := data.(float64); ok {
		return strconv.FormatFloat(intValue, 'f', -1, 64)
	} else if intValue, ok := data.(int64); ok {
		return strconv.FormatInt(intValue, 10)
	}
	return ""
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
