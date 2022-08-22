package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OptionType int32

const (
	OptionTypeClient = OptionType(iota + 1)
	OptionTypeRequest
	OptionTypeResponse
	OptionTypeLog
)

type Option struct {
	Type     OptionType `json:"type"`
	trace    *Trace
	Client   func(client *http.Client)
	Request  func(request *http.Request, trace *Trace) error
	Response func(response *http.Response, trace *Trace) error
}

/*
 * duration is time.Millisecond
 */
func WithTimeout(duration time.Duration) Option {
	return Option{
		Type: OptionTypeClient,
		Client: func(client *http.Client) {
			client.Timeout = duration
		},
	}
}

func WithContentType(typ string) Option {
	return Option{
		Type: OptionTypeRequest,
		Request: func(request *http.Request, trace *Trace) error {
			request.Header.Set("Content-Type", typ)
			if trace != nil {
				trace.Record(request, "")
			}
			return nil
		},
	}
}

func WithAuthorization(authorization string) Option {
	return Option{
		Type: OptionTypeRequest,
		Request: func(request *http.Request, trace *Trace) error {
			request.Header.Set("Authorization", authorization)
			if trace != nil {
				trace.Record(request, "")
			}
			return nil
		},
	}
}

func WithNoBody() Option {
	return Option{
		Type: OptionTypeRequest,
		Request: func(request *http.Request, trace *Trace) error {
			if trace != nil {
				trace.Record(request, "")
			}
			return nil
		},
	}
}

func WithBody(source map[string]interface{}) Option {
	return Option{
		Type: OptionTypeRequest,
		Request: func(request *http.Request, trace *Trace) error {
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			values := url.Values{}
			for k, v := range source {
				values.Add(k, getString(v))
			}
			valueString := values.Encode()
			request.Body = ioutil.NopCloser(strings.NewReader(valueString))

			if trace != nil {
				trace.Record(request, valueString)
			}

			return nil
		},
	}
}

func WithJSONBody(source interface{}) Option {
	return Option{
		Type: OptionTypeRequest,
		Request: func(request *http.Request, trace *Trace) error {
			request.Header.Set("Content-Type", "application/json")

			body, _ := json.Marshal(source)
			request.Body = ioutil.NopCloser(bytes.NewReader(body))

			if trace != nil {
				trace.Record(request, string(body))
			}

			return nil
		},
	}
}

func WithMultipartBody(source map[string]interface{}, file map[string][]byte) Option {
	return Option{
		Type: OptionTypeRequest,
		Request: func(request *http.Request, trace *Trace) error {
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
				writer.WriteField(field, getString(value))
			}

			request.Body = ioutil.NopCloser(buff)
			request.Header.Set("Content-Type", writer.FormDataContentType())

			if trace != nil {
				trace.Record(request, source, strings.Join(fileInfo, ""))
			}

			return nil
		},
	}
}

func WithBinaryBody(file []byte) Option {
	return Option{
		Type: OptionTypeRequest,
		Request: func(request *http.Request, trace *Trace) error {
			request.Body = ioutil.NopCloser(bytes.NewReader(file))
			request.Header.Set("Content-Type", "application/octet-stream")

			if trace != nil {
				trace.Record(request, fmt.Sprintf("(Binary[%d bytes])", len(file)))
			}

			return nil
		},
	}
}

func WithTrace(trace *Trace) Option {
	return Option{
		Type:  OptionTypeLog,
		trace: trace,
	}
}

func WithResponseData(data *string) Option {
	return Option{
		Type: OptionTypeResponse,
		Response: func(response *http.Response, trace *Trace) error {
			body, _ := ioutil.ReadAll(response.Body)
			bodyString := string(body)
			if data != nil {
				*data = bodyString
			}
			if trace != nil {
				trace.Record(response, bodyString)
			}
			return nil
		},
	}
}

func WithResponseJSONData(data interface{}) Option {
	return Option{
		Type: OptionTypeResponse,
		Response: func(response *http.Response, trace *Trace) error {
			body, _ := ioutil.ReadAll(response.Body)
			if trace != nil {
				trace.Record(response, string(body))
			}
			if data != nil {
				if err := json.Unmarshal(body, &data); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

type optionList []Option

func (list optionList) Client() *http.Client {
	client := &http.Client{}
	for _, option := range list {
		if option.Type != OptionTypeClient {
			continue
		}
		option.Client(client)
	}
	return client
}

func (list optionList) Request(request *http.Request) error {
	stat := false
	trace := list.Trace()
	for _, option := range list {
		if option.Type != OptionTypeRequest {
			continue
		}
		stat = true
		if err := option.Request(request, trace); err != nil {
			return err
		}
	}
	if !stat {
		WithNoBody().Request(request, trace)
	}

	return nil
}

func (list optionList) Response(response *http.Response) error {
	stat := false
	trace := list.Trace()
	for _, option := range list {
		if option.Type != OptionTypeResponse {
			continue
		}
		stat = true
		if err := option.Response(response, trace); err != nil {
			return err
		}
	}
	if !stat {
		WithResponseData(nil).Response(response, trace)
	}
	return nil
}

func (list optionList) Trace() *Trace {
	for _, option := range list {
		if option.Type != OptionTypeLog || option.trace == nil {
			continue
		}
		return option.trace
	}
	return nil
}
