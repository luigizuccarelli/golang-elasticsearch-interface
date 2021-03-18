package connectors

import (
	"net/http"

	"lmzsoftware.com/luigizuccarelli/golang-elasticsearch-interface/pkg/schema"
)

// Clients interface - the NewClientConnectors function will implement this interface
type Clients interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
	Meta(string)
	Do(req *http.Request) (*http.Response, error)
	GetDocuments(gs *schema.GenericInterface) ([]byte, error)
	ProcessDocument(method string, gs *schema.GenericInterface) error
}
