package schema

// Response schema
type Response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ElasticResponse schema
type ElasticResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string  `json:"_index"`
			Type   string  `json:"_type"`
			ID     string  `json:"_id"`
			Score  float64 `json:"_score"`
			Source struct {
				ID             string `json:"id"`
				Emailaddress   string `json:"emailaddress"`
				Firstname      string `json:"firstname"`
				Lastname       string `json:"lastname"`
				Customernumber string `json:"customerNumber"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// GenericInterface schema
// This saves us a bt of work as we don't need to update function signatures
type GenericInterface struct {
	// change the payload to point to any struct
	// in this case I'm using CustomerInfo
	Payload *CustomerInfo
	Id      string `json:"id"`
}

// CustomerInfo schema
type CustomerInfo struct {
	Id             string `json:"id,omitempty"`
	EmailAddress   string `json:"emailaddress"`
	PhoneNumber    string `json:"phonenumber,omitempty"`
	Mobile         string `json:"mobile,omitempty"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	CustomerNumber string `json:"customerNumber,omitempty"`
}

var SearchQueryTemplate = `{
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
                "query": "{{ .Payload.FirstName }}*"
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
                "query": "{{ .Payload.LastName }}*"
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
                "query": "{{ .Payload.EmailAddress }}*"
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
