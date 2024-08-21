package client

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-querystring/query"
	uuid "github.com/satori/go.uuid"
)

const (
	apiURL     = "https://api.binance.com"
	apiVersion = "v1"
)

type Client struct {
	*http.Client
	AccessKey string
	SecretKey string
}

func (c *Client) Call(method, url string, v interface{}) (interface{}, error) {
	var body []byte

	claims := claims{
		AccessKey:      c.AccessKey,
		Nonce:          uuid.NewV4(),
		StandardClaims: jwt.StandardClaims{},
	}

	values, err := query.Values(v)
	if err != nil {
		panic(err)
	}
	if len(values) > 0 {
		encodedQuery := values.Encode()

		hash := sha512.Sum512([]byte(encodedQuery))

		claims.QueryHash = hex.EncodeToString(hash[:])
		claims.QueryHashAlg = "SHA512"

		url = url + "?" + encodedQuery

		body, err = json.Marshal(values)
		if err != nil {
			return nil, err
		}
	}

	token, err := newHS256Token(c.SecretKey, claims)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, apiURL+"/"+apiVersion+url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token.Type+" "+token.SignedString)

	return getResponse(c.Client, req)
}

func getResponse(client *http.Client, req *http.Request) (interface{}, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r interface{}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	switch t := r.(type) {
	case []interface{}:
		var a []map[string]interface{}

		for _, item := range t {
			a = append(a, item.(map[string]interface{}))
		}
		r = a
	case map[string]interface{}:
		r = t
	}

	return r, nil
}
