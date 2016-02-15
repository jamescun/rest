package rest

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientNew(t *testing.T) {
	c, err := New("https://example.org", EncoderJSON{})
	if assert.NoError(t, err) {
		assert.NotNil(t, c.Encoder)
		if assert.NotNil(t, c.URL) {
			assert.Equal(t, c.URL.Scheme, "https")
			assert.Equal(t, c.URL.Host, "example.org")
		}

		if assert.NotNil(t, c.client) {
			assert.Equal(t, c.client.Timeout, 30*time.Second)
		}
	}
}

func TestClientSetTimeout(t *testing.T) {
	c := Client{client: &http.Client{Timeout: 30 * time.Second}}

	c.SetTimeout(5 * time.Second)
	assert.Equal(t, c.client.Timeout, 5*time.Second)
}

func TestClientBefore(t *testing.T) {
	c := Client{}
	c.Before(func(req *http.Request) {})
	assert.Len(t, c.reqw, 1)
}

func TestClientError(t *testing.T) {
	c := Client{}
	c.Error(func(res *http.Response, req *http.Request, dst interface{}) error { return nil })
	assert.Len(t, c.errw, 1)
}

func TestClientNewURL(t *testing.T) {
	c, _ := New("https://example.org/", nil)
	u := c.newURL("/1/foo/")
	assert.Equal(t, u.Host, "example.org")
	assert.Equal(t, u.Path, "/1/foo")

	c, _ = New("https://example.org/api/", nil)
	u = c.newURL("/foo")
	assert.Equal(t, u.Host, "example.org")
	assert.Equal(t, u.Path, "/api/foo")
}

func TestClientString(t *testing.T) {
	c, _ := New("https://example.org", nil)
	assert.Equal(t, c.String(), "Rest:https{example.org/}")
}
