// +build fake

package connectors

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/microlib/simple"
)

var count int = 0

// Connectors - overrides the real implemntation (using gocb.* dependencies)
// The file directive +build mock ensures its use (see the first line of this file)_
type Connectors struct {
	Http   *http.Client
	Logger *simple.Logger
	Flag   string
}

// Do - http wrapper
func (c *Connectors) Do(req *http.Request) (*http.Response, error) {
	if c.Flag == "error" {
		return &http.Response{}, errors.New("forced http error")
	}
	return c.Http.Do(req)
}

func (c *Connectors) Meta(id string) {
	c.Flag = id
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// NewTestConnectors - interface for mocking all tests
func NewTestConnectors(file string, code int, logger *simple.Logger) Clients {

	// we first load the json payload to simulate a call to middleware
	// for now just ignore failures.
	data, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Error(fmt.Sprintf("file data %v\n", err))
		panic(err)
	}
	httpclient := NewHttpTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: code,
			// Send response to be tested

			Body: ioutil.NopCloser(bytes.NewBufferString(string(data))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	conns := &Connectors{Http: httpclient, Logger: logger, Flag: "false"}
	return conns
}
