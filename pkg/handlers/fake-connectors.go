// +build fake

package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/microlib/simple"
	"lmzsoftware.com/luigizuccarelli/golang-elasticsearch-interface/pkg/connectors"
	"lmzsoftware.com/luigizuccarelli/golang-elasticsearch-interface/pkg/schema"
)

var count int = 0

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Inject (force) readAll test error")
}

type FakeConnectors struct {
	Http   *http.Client
	Logger *simple.Logger
	Name   string
	Force  string
}

// GetAllDocument - simulate (fake) the call
func (c *FakeConnectors) GetDocuments(gs *schema.GenericInterface) ([]byte, error) {
	if c.Force == "error" {
		return []byte(""), errors.New("GetAllDocument forced error")
	}
	d, _ := ioutil.ReadFile("../../tests/all-data.json")
	return d, nil
}

// ProcessDocument - simulate (fake) the call
func (c *FakeConnectors) ProcessDocument(method string, gs *schema.GenericInterface) error {
	if c.Force == "error" {
		return errors.New("ProcessDocument forced error")
	}
	return nil
}

// logger wrapper
func (r *FakeConnectors) Error(msg string, val ...interface{}) {
	r.Logger.Error(fmt.Sprintf(msg, val...))
}

func (r *FakeConnectors) Info(msg string, val ...interface{}) {
	r.Logger.Info(fmt.Sprintf(msg, val...))
}

func (r *FakeConnectors) Debug(msg string, val ...interface{}) {
	r.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (r *FakeConnectors) Trace(msg string, val ...interface{}) {
	r.Logger.Trace(fmt.Sprintf(msg, val...))
}

func (c *FakeConnectors) Meta(flag string) {
	c.Force = flag
}

// Do - http wrapper
func (c *FakeConnectors) Do(req *http.Request) (*http.Response, error) {
	return c.Http.Do(req)
}

// NewTestConnectors - inject our test connectors
func NewTestConnectors(filename string, code int, logger *simple.Logger) connectors.Clients {

	// we first load the json payload to simulate response data
	// for now just ignore failures.
	file, _ := ioutil.ReadFile(filename)
	logger.Trace(fmt.Sprintf("File %s with data %s", filename, string(file)))

	conn := &FakeConnectors{Logger: logger, Name: "test"}
	return conn
}
