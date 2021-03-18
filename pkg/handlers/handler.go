package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"lmzsoftware.com/luigizuccarelli/golang-elasticsearch-interface/pkg/connectors"
	"lmzsoftware.com/luigizuccarelli/golang-elasticsearch-interface/pkg/schema"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
	ERROR           string = "ERROR"
)

// DocumentHandler - all data api function handler
func DocumentHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var gs *schema.GenericInterface

	addHeaders(w, r)
	//vars := mux.Vars(r)

	con.Info("DocumentHandler method %s", r.Method)

	// ensure we don't have nil - it will cause a null pointer exception
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "DocumentHandler body data error :  %v"
		con.Error(msg, err)
		response := responseFormat(http.StatusInternalServerError, "ERROR", w, msg, err)
		fmt.Fprintf(w, "%s", response)
		return
	}

	con.Trace("DocumentHandler request body : %s", string(body))

	// unmarshal result - ensures correct json struct  (validation)
	errs := json.Unmarshal(body, &gs)
	if errs != nil {
		msg := "DocumentHandler could not unmarshal input data to schema %v"
		con.Error(msg, errs)
		response := responseFormat(http.StatusInternalServerError, "ERROR", w, msg, errs)
		fmt.Fprintf(w, "%s", response)
		return
	}
	e := con.ProcessDocument(r.Method, gs)
	if e != nil {
		msg := "DocumentHandler could not insert(elasticsearch) %v"
		con.Error(msg, e)
		response := responseFormat(http.StatusInternalServerError, "ERROR", w, msg, e)
		fmt.Fprintf(w, "%s", response)
		return
	}

	response := responseFormat(http.StatusOK, "OK", w, "DcoumentHandler call successfull")
	fmt.Fprintf(w, "%s", response)
	return
}

// SearchHandler - all data api function handler
func SearchHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var gs *schema.GenericInterface

	addHeaders(w, r)
	//vars := mux.Vars(r)

	con.Info("SearchHandler method %s", r.Method)

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := "SearchHandler body data error :  %v"
			con.Error(msg, err)
			response := responseFormat(http.StatusInternalServerError, "ERROR", w, msg, err)
			fmt.Fprintf(w, "%s\n", response)
			return
		}

		con.Trace("SearchHandler request body : %s", string(body))

		// unmarshal result - ensures correct json struct  (validation)
		errs := json.Unmarshal(body, &gs)
		if errs != nil {
			msg := "SearchHandler could not unmarshal input data to schema %v"
			con.Error(msg, errs)
			response := responseFormat(http.StatusInternalServerError, "ERROR", w, msg, errs)
			fmt.Fprintf(w, "%s\n", response)
			return
		}

		// ensure parameter is valid
		es, err := con.GetDocuments(gs)
		if err != nil {
			msg := "SearchHandler (elastisearch) lookup %v"
			con.Error(msg, err)
			response := responseFormat(http.StatusInternalServerError, "ERROR", w, msg, errs)
			fmt.Fprintf(w, "%s\n", response)
			return
		}

		var eschema *schema.ElasticResponse
		// we only want the hits array
		json.Unmarshal(es, &eschema)
		b, _ := json.MarshalIndent(eschema.Hits.Hits, "", "	")
		response := responseFormat(http.StatusOK, "OK", w, "%s", string(b))
		fmt.Fprintf(w, "%s\n", response)
		return
	}
}

// utility functions

// responsFormat
func responseFormat(code int, status string, w http.ResponseWriter, msg string, val ...interface{}) string {
	response := `{"Code":"` + strconv.Itoa(code) + `", "Status": "` + status + `", "Message":"` + fmt.Sprintf(msg, val...) + `"}`
	w.WriteHeader(code)
	return response
}

// IsAlive - used for readiness and liveness probes
func IsAlive(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \""+os.Getenv("NAME")+"\" }")
	return
}

// addHeaders - headers (with cors)
func addHeaders(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("API-KEY") != "" {
		w.Header().Set("API_KEY_PT", r.Header.Get("API_KEY"))
	}
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept-Language")
}
