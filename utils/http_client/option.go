package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/1t-data-kit/go-kit/base"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ClientWrapper func(client *http.Client, trace Trace)
type RequestWrapper func(request *http.Request, trace Trace) error
type ResponseWrapper func(response *http.Response, trace Trace) error

func WithTimeout(timeoutMillisecond int64) base.Option {
	var wrapper ClientWrapper = func(client *http.Client, trace Trace) {
		client.Timeout = time.Duration(timeoutMillisecond) * time.Millisecond
	}
	return base.NewOption(wrapper)
}

func WithContentType(typ string) base.Option {
	var wrapper RequestWrapper = func(request *http.Request, trace Trace) error {
		request.Header.Set("Content-Type", typ)
		if !traceIsNil(trace) {
			trace.SetUrl(request.URL.String())
			trace.SetHeader(base.AssertString(request.Header))
		}
		return nil
	}
	return base.NewOption(wrapper)
}

func WithAuthorization(authorization string) base.Option {
	var wrapper RequestWrapper = func(request *http.Request, trace Trace) error {
		request.Header.Set("Authorization", authorization)
		if !traceIsNil(trace) {
			trace.SetUrl(request.URL.String())
			trace.SetHeader(base.AssertString(request.Header))
		}
		return nil
	}
	return base.NewOption(wrapper)
}

func WithEmptyBody() base.Option {
	var wrapper RequestWrapper = func(request *http.Request, trace Trace) error {
		if !traceIsNil(trace) {
			trace.SetUrl(request.URL.String())
			trace.SetHeader(base.AssertString(request.Header))
		}
		return nil
	}
	return base.NewOption(wrapper)
}

func WithQuery(query string) base.Option {
	var wrapper RequestWrapper = func(request *http.Request, trace Trace) error {
		values := url.Values{}
		for _, kv := range strings.Split(query, "&") {
			k := kv
			v := ""
			sep := strings.Index(kv, "=")
			if sep > -1 {
				k = kv[0:sep]
				v = kv[sep+1:]
			}
			values.Add(k, v)
		}
		request.URL.RawQuery = values.Encode()

		if !traceIsNil(trace) {
			trace.SetUrl(request.URL.String())
			trace.SetHeader(base.AssertString(request.Header))
		}

		return nil
	}
	return base.NewOption(wrapper)
}

func WithBody(source map[string]interface{}) base.Option {
	var wrapper RequestWrapper = func(request *http.Request, trace Trace) error {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		values := url.Values{}
		for k, v := range source {
			values.Add(k, base.AssertString(v))
		}
		valueString := values.Encode()
		request.Body = io.NopCloser(strings.NewReader(valueString))

		if !traceIsNil(trace) {
			trace.SetUrl(request.URL.String())
			trace.SetHeader(base.AssertString(request.Header))
			trace.SetRequest(valueString)
		}

		return nil
	}
	return base.NewOption(wrapper)
}

func WithJSONBody(source interface{}, opt ...base.Option) base.Option {
	var wrapper RequestWrapper = func(request *http.Request, trace Trace) error {
		mustEscape := false
		if opts := base.Options(opt).Filter(func(item base.Option) bool {
			if _, ok := item.Value().(bool); ok {
				return true
			}
			return false
		}); len(opts) > 0 {
			mustEscape = opts[0].Value().(bool)
		}
		request.Header.Set("Content-Type", "application/json")

		buf := bytes.NewBuffer([]byte{})
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(mustEscape)
		if err := encoder.Encode(source); err != nil {
			return err
		}

		body := buf.Bytes()
		request.Body = io.NopCloser(bytes.NewReader(body))

		if !traceIsNil(trace) {
			trace.SetUrl(request.URL.String())
			trace.SetHeader(base.AssertString(request.Header))
			trace.SetRequest(string(body))
		}

		return nil
	}
	return base.NewOption(wrapper)
}

func WithMultipartBody(source map[string]interface{}, file map[string][]byte) base.Option {
	var wrapper RequestWrapper = func(request *http.Request, trace Trace) error {
		buff := &bytes.Buffer{}
		writer := multipart.NewWriter(buff)
		defer writer.Close()

		fileInfo := make([]string, len(file))
		for name, content := range file {
			fileWriter, err := writer.CreateFormFile(name, name)
			if err != nil {
				return err
			}
			if _, err = fileWriter.Write(content); err != nil {
				return err
			}
			fileInfo = append(fileInfo, fmt.Sprintf("(Binary[%s][%d bytes])", name, len(content)))
		}
		for field, value := range source {
			writer.WriteField(field, base.AssertString(value))
		}

		request.Body = io.NopCloser(buff)
		request.Header.Set("Content-Type", writer.FormDataContentType())

		if !traceIsNil(trace) {
			trace.SetUrl(request.URL.String())
			trace.SetHeader(base.AssertString(request.Header))
			trace.SetRequest(base.AssertStrings(",", source, strings.Join(fileInfo, "")))
		}

		return nil
	}
	return base.NewOption(wrapper)
}

func WithBinaryBody(file []byte) base.Option {
	var wrapper RequestWrapper = func(request *http.Request, trace Trace) error {
		request.Body = io.NopCloser(bytes.NewReader(file))
		request.Header.Set("Content-Type", "application/octet-stream")

		if !traceIsNil(trace) {
			trace.SetUrl(request.URL.String())
			trace.SetHeader(base.AssertString(request.Header))
			trace.SetRequest(fmt.Sprintf("(Binary[%d bytes])", len(file)))
		}

		return nil
	}
	return base.NewOption(wrapper)
}

func WithData(data *string) base.Option {
	var wrapper ResponseWrapper = func(response *http.Response, trace Trace) error {
		body, _ := io.ReadAll(response.Body)
		bodyString := string(body)
		if data != nil {
			*data = bodyString
		}
		if !traceIsNil(trace) {
			trace.SetHttpCode(response.StatusCode)
			trace.SetResponse(bodyString)
		}
		return nil
	}
	return base.NewOption(wrapper)
}

func WithJsonData(data interface{}) base.Option {
	var wrapper ResponseWrapper = func(response *http.Response, trace Trace) error {
		body, _ := io.ReadAll(response.Body)
		if !traceIsNil(trace) {
			trace.SetHttpCode(response.StatusCode)
			trace.SetResponse(string(body))
		}
		if data != nil {
			if err := json.Unmarshal(body, &data); err != nil {
				return err
			}
		}

		return nil
	}
	return base.NewOption(wrapper)
}

func WithTrace(trace Trace) base.Option {
	return base.NewOption(trace)
}
