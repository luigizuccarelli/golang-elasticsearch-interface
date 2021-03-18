package connectors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/segmentio/ksuid"
	"lmzsoftware.com/luigizuccarelli/golang-elasticsearch-interface/pkg/schema"
)

func (c *Connectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

// ProcessDocument - function that inserts data to elasticsearch
func (c *Connectors) ProcessDocument(method string, gs *schema.GenericInterface) error {
	var id string
	// this is redundant but needs to be done
	c.Info("Function ProcessDocument input data %v", gs.Payload)

	if gs.Payload == nil {
		c.Error("ProcessDocument payload is empty")
		return errors.New("payload is empty")
	}
	doc, err := json.MarshalIndent(gs.Payload, "", "  ")
	if err != nil {
		c.Error("Function ProcessDocument marshalling data from struct: %v", err)
		return err
	}
	if gs.Id == "" {
		id = ksuid.New().String()
	} else {
		id = gs.Id
	}
	req, err := http.NewRequest(method, os.Getenv("ELASTICSEARCH_URL")+"/"+os.Getenv("INDEX")+"/_doc/"+id, bytes.NewBuffer(doc))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		c.Error("Function ProcessDocument request %v", err)
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Error("Function ProcessDocument reading response %v", err)
		return err
	}

	c.Debug("Function ProcessDocument result %v", string(data))
	c.Info("Function ProcessDocument returned successfully")
	return nil
}

// GetDocument - use search firstname,lastname, email
func (c *Connectors) GetDocuments(gs *schema.GenericInterface) ([]byte, error) {

	query := `{
   "from" : 0,
	 "size": 10,
	 "query": {
    "bool": {
      "must": [],
      "filter": [
        {
          "bool": {
            "should": [
              {
                "query_string": {
                  "fields": [
                    "firstname"
                  ],
                  "query": "` + gs.Payload.FirstName + `*"
                }
              }
            ],
            "minimum_should_match": 1
          }
        },
				{
          "bool": {
            "should": [
              {
                "query_string": {
                  "fields": [
                    "lastname"
                  ],
                  "query": "` + gs.Payload.LastName + `*"
                }
              }
            ],
            "minimum_should_match": 1
          }
        },
				{
          "bool": {
            "should": [
              {
                "query_string": {
                  "fields": [
                    "emailaddress"
                  ],
                  "query": "` + gs.Payload.EmailAddress + `*"
                }
              }
            ],
            "minimum_should_match": 1
          }
        }
      ],
      "should": [],
      "must_not": []
    }
   }
 }`

	c.Info("Query string %s", query)
	req, err := http.NewRequest("POST", os.Getenv("ELASTICSEARCH_URL")+"/"+os.Getenv("INDEX")+"/_search", bytes.NewBuffer([]byte(query)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		c.Error("Function GetDocument request %v", err)
		return []byte(""), err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Error("Function GetDocument reading response %v", err)
		return []byte(""), err
	}

	c.Debug("Function GetDocument result %v", string(data))
	c.Info("Function GetDocument returned successfully")
	return data, nil
}
