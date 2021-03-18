// +build fake

package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/microlib/simple"
)

// Makes use of the fake-connectors.go (in this package) for testing via the +build fake directive

func TestHandlers(t *testing.T) {

	var logger = &simple.Logger{Level: "info"}

	t.Run("IsAlive : should pass", func(t *testing.T) {
		var STATUS int = 200

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v2/sys/info/isalive", nil)
		NewTestConnectors("../../tests/emails.json", STATUS, logger)
		handler := http.HandlerFunc(IsAlive)
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "IsAlive", rr.Code, STATUS))
		}
	})

	t.Run("DoumentHandler : should pass (put)", func(t *testing.T) {
		var STATUS int = 200

		requestPayload, _ := ioutil.ReadFile("../../tests/generic-interface.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/customers", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/insert.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DocumentHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "DocumentHandler ", rr.Code, STATUS))
		}
	})

	t.Run("DoumentHandler : should fail (force readall error)", func(t *testing.T) {
		var STATUS int = 500

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/customers", errReader(0))
		conn := NewTestConnectors("../../tests/insert.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DocumentHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "DocumentHandler ", rr.Code, STATUS))
		}
	})

	t.Run("DoumentHandler : should fail (force json input)", func(t *testing.T) {
		var STATUS int = 500

		requestPayload, _ := ioutil.ReadFile("../../tests/generic-interface-error.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/customers", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/insert.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DocumentHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "DocumentHandler ", rr.Code, STATUS))
		}
	})

	t.Run("DoumentHandler : should fail (input json empty)", func(t *testing.T) {
		var STATUS int = 500

		requestPayload, _ := ioutil.ReadFile("../../tests/generic-interface-empty.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/customers", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
		conn.Meta("error")
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DocumentHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "DocumentHandler ", rr.Code, STATUS))
		}
	})

	t.Run("SearchHandler : should pass (post)", func(t *testing.T) {
		var STATUS int = 200

		requestPayload, _ := ioutil.ReadFile("../../tests/generic-interface.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/search/customers", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SearchHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "DocumentHandler ", rr.Code, STATUS))
		}
	})

	t.Run("SearchHandler : should fail (force readall error)", func(t *testing.T) {
		var STATUS int = 500

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/search/customers", errReader(0))
		conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SearchHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "DocumentHandler ", rr.Code, STATUS))
		}
	})

	t.Run("SearchHandler : should fail (force json error)", func(t *testing.T) {
		var STATUS int = 500

		requestPayload, _ := ioutil.ReadFile("../../tests/generic-interface-error.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/search/customers", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SearchHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "DocumentHandler ", rr.Code, STATUS))
		}
	})

	t.Run("SearchHandler : should fail (force GetDocuments error)", func(t *testing.T) {
		var STATUS int = 500

		requestPayload, _ := ioutil.ReadFile("../../tests/generic-interface.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/search/customers", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
		conn.Meta("error")
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SearchHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "DocumentHandler ", rr.Code, STATUS))
		}
	})

	/*
		t.Run("CustomerHandler : should pass (get)", func(t *testing.T) {
			var STATUS int = 200

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/customers/12321312", nil)
			conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				CustomerHandler(w, r, conn)
			})
			//Hack to try to fake gorilla/mux vars
			vars := map[string]string{
				"id": "232131313",
			}
			req = mux.SetURLVars(req, vars)
			handler.ServeHTTP(rr, req)
			body, e := ioutil.ReadAll(rr.Body)
			if e != nil {
				t.Fatalf("Should not fail : found error %v", e)
			}
			logger.Trace(fmt.Sprintf("Response %s", string(body)))
			// ignore errors here
			if rr.Code != STATUS {
				t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CustomerHandler ", rr.Code, STATUS))
			}
		})

		t.Run("CustomerHandler : should fail (get)", func(t *testing.T) {
			var STATUS int = 500

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/customers/12321312", nil)
			conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				CustomerHandler(w, r, conn)
			})
			// no mux vars set
			handler.ServeHTTP(rr, req)
			body, e := ioutil.ReadAll(rr.Body)
			if e != nil {
				t.Fatalf("Should not fail : found error %v", e)
			}
			logger.Trace(fmt.Sprintf("Response %s", string(body)))
			// ignore errors here
			if rr.Code != STATUS {
				t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CustomerHandler ", rr.Code, STATUS))
			}
		})

		t.Run("CustomerHandler : should fail (get forced couchbase error)", func(t *testing.T) {
			var STATUS int = 500

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/customers/12321312", nil)
			conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
			conn.Meta("true")
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				CustomerHandler(w, r, conn)
			})
			//Hack to try to fake gorilla/mux vars
			vars := map[string]string{
				"customernumber": "232131313",
			}
			req = mux.SetURLVars(req, vars)
			handler.ServeHTTP(rr, req)
			body, e := ioutil.ReadAll(rr.Body)
			if e != nil {
				t.Fatalf("Should not fail : found error %v", e)
			}
			logger.Trace(fmt.Sprintf("Response %s", string(body)))
			// ignore errors here
			if rr.Code != STATUS {
				t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CustomerHandler ", rr.Code, STATUS))
			}
		})

		t.Run("CustomerHandler : should fail (forced read error)", func(t *testing.T) {
			var STATUS int = 500

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/customers/12321312", errReader(0))
			conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				CustomerHandler(w, r, conn)
			})
			handler.ServeHTTP(rr, req)
			body, e := ioutil.ReadAll(rr.Body)
			if e != nil {
				t.Fatalf("Should not fail : found error %v", e)
			}
			logger.Trace(fmt.Sprintf("Response %s", string(body)))
			// ignore errors here
			if rr.Code != STATUS {
				t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CustomerHandler ", rr.Code, STATUS))
			}
		})

		t.Run("CustomerHandler : should fail (json request data)", func(t *testing.T) {
			var STATUS int = 500

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/customers/12321312", bytes.NewBuffer([]byte("{ accounts")))
			conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				CustomerHandler(w, r, conn)
			})
			handler.ServeHTTP(rr, req)
			body, e := ioutil.ReadAll(rr.Body)
			if e != nil {
				t.Fatalf("Should not fail : found error %v", e)
			}
			logger.Trace(fmt.Sprintf("Response %s", string(body)))
			// ignore errors here
			if rr.Code != STATUS {
				t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CustomerHandler ", rr.Code, STATUS))
			}
		})

		t.Run("CustomerHandler : should fail (WriteDocument)", func(t *testing.T) {
			var STATUS int = 500

			requestPayload, _ := ioutil.ReadFile("../../tests/all-data.json")
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/customers/12321312", bytes.NewBuffer(requestPayload))
			conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				CustomerHandler(w, r, conn)
			})
			conn.Meta("true")
			handler.ServeHTTP(rr, req)
			body, e := ioutil.ReadAll(rr.Body)
			if e != nil {
				t.Fatalf("Should not fail : found error %v", e)
			}
			logger.Trace(fmt.Sprintf("Response %s", string(body)))
			// ignore errors here
			if rr.Code != STATUS {
				t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "CustomerHandler ", rr.Code, STATUS))
			}
		})

		t.Run("SearchCustomerHandler : should pass (post)", func(t *testing.T) {
			var STATUS int = 200

			requestPayload := `{ "lastName":"test","firstName":"test"}`
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/search/customer", bytes.NewBuffer([]byte(requestPayload)))
			conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				SearchCustomerHandler(w, r, conn)
			})
			handler.ServeHTTP(rr, req)
			body, e := ioutil.ReadAll(rr.Body)
			if e != nil {
				t.Fatalf("Should not fail : found error %v", e)
			}
			logger.Trace(fmt.Sprintf("Response %s", string(body)))
			// ignore errors here
			if rr.Code != STATUS {
				t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "SearchCustomerHandler ", rr.Code, STATUS))
			}
		})

		t.Run("SearchCustomerHandler : should fail (forced read error)", func(t *testing.T) {
			var STATUS int = 500

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/search/customer", errReader(0))
			conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				SearchCustomerHandler(w, r, conn)
			})
			handler.ServeHTTP(rr, req)
			body, e := ioutil.ReadAll(rr.Body)
			if e != nil {
				t.Fatalf("Should not fail : found error %v", e)
			}
			logger.Trace(fmt.Sprintf("Response %s", string(body)))
			// ignore errors here
			if rr.Code != STATUS {
				t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "SearchCustomerHandler ", rr.Code, STATUS))
			}
		})

		t.Run("AccountsHandler : should fail (json request data)", func(t *testing.T) {
			var STATUS int = 500

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer([]byte("{ firstName")))
			conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				SearchCustomerHandler(w, r, conn)
			})
			handler.ServeHTTP(rr, req)
			body, e := ioutil.ReadAll(rr.Body)
			if e != nil {
				t.Fatalf("Should not fail : found error %v", e)
			}
			logger.Trace(fmt.Sprintf("Response %s", string(body)))
			// ignore errors here
			if rr.Code != STATUS {
				t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "SearchCustomerHandler ", rr.Code, STATUS))
			}
		})

		t.Run("SearchCustomerHandler : should fail (GetDocuments)", func(t *testing.T) {
			var STATUS int = 500

			requestPayload := `{ "firstName":"test","lastName":"test"}`
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/search/customers", bytes.NewBuffer([]byte(requestPayload)))
			conn := NewTestConnectors("../../tests/all-data.json", STATUS, logger)
			conn.Meta("true")
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				SearchCustomerHandler(w, r, conn)
			})
			handler.ServeHTTP(rr, req)
			body, e := ioutil.ReadAll(rr.Body)
			if e != nil {
				t.Fatalf("Should not fail : found error %v", e)
			}
			logger.Trace(fmt.Sprintf("Response %s", string(body)))
			// ignore errors here
			if rr.Code != STATUS {
				t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "AccountsHandler ", rr.Code, STATUS))
			}
		})
	*/
}
