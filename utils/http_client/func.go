package http_client

import (
	"context"
	"github.com/1t-data-kit/go-kit/base"
	"net/http"
)

func GET(ctx context.Context, url string, option ...base.Option) error {
	return Do(ctx, "GET", url, option...)
}

func POST(ctx context.Context, url string, source map[string]interface{}, option ...base.Option) error {
	return Do(ctx, "POST", url, append(option, WithBody(source))...)
}

func POSTJson(ctx context.Context, url string, source interface{}, option ...base.Option) error {
	return Do(ctx, "POST", url, append(option, WithJSONBody(source))...)
}

func POSTBinary(ctx context.Context, url string, file []byte, option ...base.Option) error {
	return Do(ctx, "POST", url, append(option, WithBinaryBody(file))...)
}

func POSTMultipart(ctx context.Context, url string, source map[string]interface{}, file map[string][]byte, option ...base.Option) error {
	return Do(ctx, "POST", url, append(option, WithMultipartBody(source, file))...)
}

func Do(ctx context.Context, method string, url string, option ...base.Option) error {
	client := &http.Client{}
	_trace := trace(option...)
	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return err
	}
	wrapClient(client, _trace, option...)
	if err = wrapRequest(request, _trace, option...); err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if err = wrapResponse(response, _trace, option...); err != nil {
		return err
	}
	return nil
}

func trace(option ...base.Option) Trace {
	if len(option) == 0 {
		return nil
	}
	if traces := base.Options(option).Filter(func(item base.Option) bool {
		if _, ok := item.Value().(Trace); ok {
			return true
		}
		return false
	}); len(traces) > 0 {
		return traces[len(traces)-1].Value().(Trace)
	}
	return nil
}

func wrapClient(client *http.Client, trace Trace, option ...base.Option) {
	if len(option) == 0 {
		return
	}

	if clientWrapperOptions := base.Options(option).Filter(func(item base.Option) bool {
		if _, ok := item.Value().(func(client *http.Client, trace Trace)); ok {
			return true
		}
		return false
	}); len(clientWrapperOptions) > 0 {
		for _, clientWrapperOption := range clientWrapperOptions {
			clientWrapperOption.Value().(func(client *http.Client, trace Trace))(client, trace)
		}
	}
}

func wrapRequest(request *http.Request, trace Trace, option ...base.Option) error {
	if len(option) == 0 {
		return nil
	}

	requestWrapperOptions := base.Options(option).Filter(func(item base.Option) bool {
		if _, ok := item.Value().(RequestWrapper); ok {
			return true
		}
		return false
	})
	if len(requestWrapperOptions) == 0 {
		requestWrapperOptions = append(requestWrapperOptions, WithEmptyBody())
	}

	errors := base.NewErrors()
	for _, requestWrapperOption := range requestWrapperOptions {
		if err := requestWrapperOption.Value().(RequestWrapper)(request, trace); err != nil {
			errors.Append(err)
		}
	}
	return errors.Error()
}

func wrapResponse(response *http.Response, trace Trace, option ...base.Option) error {
	if len(option) == 0 {
		return nil
	}

	responseWrapperOptions := base.Options(option).Filter(func(item base.Option) bool {
		if _, ok := item.Value().(ResponseWrapper); ok {
			return true
		}
		return false
	})
	if len(responseWrapperOptions) == 0 {
		responseWrapperOptions = append(responseWrapperOptions, WithData(nil))
	}

	errors := base.NewErrors()
	for _, responseWrapperOption := range responseWrapperOptions {
		if err := responseWrapperOption.Value().(ResponseWrapper)(response, trace); err != nil {
			errors.Append(err)
		}
	}
	return errors.Error()
}
