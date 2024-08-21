package client

import (
	"net/http"

	"github.com/google/go-querystring/query"
)

type QuotationClient struct {
	*http.Client
}

func (qc *QuotationClient) Call(url string, v interface{}) (interface{}, error) {
	values, err := query.Values(v)
	if err != nil {
		panic(err)
	}
	encodedQuery := values.Encode()

	req, err := http.NewRequest("GET", apiURL+"/"+apiVersion+url+"?"+encodedQuery, nil)
	if err != nil {
		return nil, err
	}

	return getResponse(qc.Client, req)
}
