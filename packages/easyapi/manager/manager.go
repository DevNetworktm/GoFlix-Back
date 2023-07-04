package manager

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Params        map[string]string
	Query         url.Values
	Body          io.Reader
	Form          url.Values
	Method        string
	Header        http.Header
	ContentLength int64
	Host          string
	RequestHTTP   *http.Request
	customVar     map[string]interface{}
}

func (req *Request) GetHeader(name string) string {
	value := req.Header.Get(name)
	return value
}

func (req *Request) GetRequestVar(name string) interface{} {
	value := req.customVar[name]
	return value
}

func (req *Request) SetRequestVar(content map[string]interface{}) {
	req.customVar = content
}

// Response

type Response struct {
	Write    http.ResponseWriter
	status   int16
	jsonDate interface{}
}

func (res *Response) Json(data interface{}) {
	res.Write.Header().Add("Content-Type", "application/json")
	res.Write.WriteHeader(int(res.status))
	json.NewEncoder(res.Write).Encode(data)
}

func (res *Response) SendStatus(status int16) {
	res.Write.WriteHeader(int(status))
}

func (res *Response) Send(message string) {
	res.Write.WriteHeader(int(res.status))
	res.Write.Write([]byte(message))
}

func (res *Response) Status(status int16) *Response {
	res.status = status
	return res
}
