// +build fake

package connectors

import (
	"fmt"
	"testing"

	"github.com/microlib/simple"
	"lmzsoftware.com/luigizuccarelli/golang-elasticsearch-interface/pkg/schema"
)

// For testing we plugin (make use of fake-connectors.go) via the +build fake directive (see first line of this file)

func TestConnections(t *testing.T) {

	var logger = &simple.Logger{Level: "trace"}

	t.Run("Logging : should pass", func(t *testing.T) {
		con := NewTestConnectors("../../tests/insert.json", 200, logger)
		con.Info("Log Info")
		con.Debug("Log Debug")
		con.Trace("Log Trace")
		con.Error("Log Error")
	})

	t.Run("ProcessDocument : should pass", func(t *testing.T) {
		con := NewTestConnectors("../../tests/insert.json", 200, logger)
		ci := &schema.CustomerInfo{LastName: "fritz"}
		gs := &schema.GenericInterface{Id: "123456", Payload: ci}
		err := con.ProcessDocument("PUT", gs)
		if err != nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error should be nil) -  got (%v) wanted (%v)", "ProcessDocument", err, nil))
		}
	})

	t.Run("ProcessDocument : should fail", func(t *testing.T) {
		con := NewTestConnectors("../../tests/insert.json", 500, logger)
		con.Meta("error")
		ci := &schema.CustomerInfo{LastName: ""}
		gs := &schema.GenericInterface{Payload: ci}
		err := con.ProcessDocument("PUT", gs)
		if err == nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error should not be nil) -  got (%v) wanted (%v)", "ProcessDocument", nil, "error"))
		}
	})

	t.Run("ProcessDocument : should fail", func(t *testing.T) {
		con := NewTestConnectors("../../tests/insert.json", 500, logger)
		gs := &schema.GenericInterface{}
		err := con.ProcessDocument("PUT", gs)
		if err == nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error shouldi not be nil) -  got (%v) wanted (%v)", "ProcessDocument", nil, "error"))
		}
	})

	t.Run("GetDocuments : should pass", func(t *testing.T) {
		con := NewTestConnectors("../../tests/all-data.json", 200, logger)
		// FirstName,LastName,EmailAddress are mandatory
		ci := &schema.CustomerInfo{LastName: "fritz", FirstName: "", EmailAddress: ""}
		gs := &schema.GenericInterface{Id: "123456", Payload: ci}
		data, err := con.GetDocuments(gs)
		if err != nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error should be nil) -  got (%v) wanted (%v)", "GetDocuments", err, nil))
		}
		con.Info("Response data %s", string(data))
	})
}
