package rest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// Request the user function for executing API requests against the client.
// it only requires a destination type to unmarshal the result (or nil if
// expecting an empty response). additionally it may take a object to marshal
// as the request body and any number of query paramaters.
type Request func(dst interface{}, src ...interface{}) error

// Requestware is a handler executed before any client request for the purpose
// of modifying the request before being set (i.e. setting headers).
type Requestware func(req *http.Request)

// Errorware is a handler executed after a client request has completed and
// the remote host returned a non-2xx status code. it is given the upstream
// response, original request, and the object expected to unmarshal to.
type Errorware func(res *http.Response, req *http.Request, dst interface{}) error

// Encoder is the interface implemented by serialization libraries/wrappers
// to handle encoding the request body and decoding and response body.
// additionally an encoder is expected to set an appropriate Content-Type
// header.
type Encoder interface {
	Encode(req *http.Request, src interface{}) ([]byte, error)
	Decode(res *http.Response, dst interface{}) error
}

// QueryParam is the interface implemented by objects to add or modify
// query parameters before the execution of a request. it takes the current
// query parameters and is expected to return them with its changes.
type QueryParam interface {
	Value(url.Values) url.Values
}

// Client represents a remote host exposing an HTTP(S) API.
type Client struct {
	// base URL (endpoint of API) for all requests
	URL *url.URL

	// request/response body encoder/decoder
	Encoder Encoder

	reqw []Requestware
	errw []Errorware

	client *http.Client
}

// return a new REST API client from a url.
func New(urlStr string, enc Encoder) (*Client, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return &Client{
		URL:     u,
		Encoder: enc,

		client: &http.Client{
			// 30 second default timeout
			Timeout: 30 * time.Second,
		},
	}, nil
}

// set the default request timeout for client.
func (c *Client) SetTimeout(d time.Duration) {
	c.client.Timeout = d
}

// add Requestware fn to be executed before every request.
func (c *Client) Before(fn Requestware) {
	c.reqw = append(c.reqw, fn)
}

// add Errorware fn to be executed after every unsuccessful request.
func (c *Client) Error(fn Errorware) {
	c.errw = append(c.errw, fn)
}

// return a new Request function for the method and url to execute against
// the client.
func (c *Client) Request(method, urlStr string) Request {
	// TODO: "make it functional, then make it beautiful, then make it performant."
	// - Erlang and OTP in Action, Joe Armstrong.
	return func(dst interface{}, src ...interface{}) (err error) {
		u := c.newURL(urlStr)
		req := &http.Request{
			Method:     method,
			URL:        u,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			Host:       u.Host,
		}

		// build query parameters and body
		v := make(url.Values)

		for i := 0; i < len(src); i++ {
			if q, ok := src[i].(QueryParam); ok {
				v = q.Value(v)
			} else if req.Body != nil {
				switch b := src[i].(type) {
				case string:
					req.Body = ioutil.NopCloser(strings.NewReader(b))
					req.ContentLength = int64(len(b))

				case []byte:
					req.Body = ioutil.NopCloser(bytes.NewReader(b))
					req.ContentLength = int64(len(b))

				case io.ReadCloser:
					req.Body = b

				case io.Reader:
					req.Body = ioutil.NopCloser(b)

				default:
					if c.Encoder != nil {
						var data []byte
						data, err = c.Encoder.Encode(req, src[i])
						if err != nil {
							return
						}

						req.Body = ioutil.NopCloser(bytes.NewReader(data))
						req.ContentLength = int64(len(data))
					}
				}
			}
		}

		req.URL.RawQuery = v.Encode()

		// execute requestware
		for i := 0; i < len(c.reqw); i++ {
			c.reqw[i](req)
		}

		// execute request
		res, err := c.client.Do(req)
		if err != nil {
			return
		}
		defer res.Body.Close()

		// non-2xx status code, execute errorware
		if res.StatusCode < 200 || res.StatusCode > 299 {
			for i := 0; i < len(c.errw); i++ {
				err = c.errw[i](res, req, dst)
				if err != nil {
					return
				}
			}

			return fmt.Errorf("http status %d", res.StatusCode)
		}

		if c.Encoder != nil {
			err = c.Encoder.Decode(res, dst)
			if err != nil {
				return
			}
		}

		return
	}
}

// return copy of client url for API endpoint
func (c *Client) newURL(urlStr string) *url.URL {
	return &url.URL{
		Scheme: c.URL.Scheme,
		Opaque: c.URL.Opaque,
		Host:   c.URL.Host,
		Path:   path.Join(c.URL.Path, urlStr),
	}
}

// return string representation of client.
func (c *Client) String() string {
	return fmt.Sprintf("Rest:%s{%s/%s}", c.URL.Scheme, c.URL.Host, c.URL.Path)
}
