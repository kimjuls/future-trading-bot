package binanceapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

const baseUrl = "https://fapi.binance.com"

type Client struct {
	*http.Client
	AccessKey string
	SecretKey string
}

func (c *Client) Request(method string, path string, params interface{}) (interface{}, error) {
	url := fmt.Sprintf("%s%s", baseUrl, path)
	var reqBody []byte

	values, err := query.Values(params)
	if err != nil {
		panic(err)
	}
	now := time.Now()
	values.Add("timestamp", fmt.Sprintf("%d", now.UnixMilli()))
	if len(values) > 0 {
		encoded := values.Encode()
		signature := sign(&encoded, &c.SecretKey)
		values.Add("signature", signature)
		url = fmt.Sprintf("%s?%s", url, values.Encode())
	}
	reqBody, err = json.Marshal(values)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.AccessKey)
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var r interface{}
	err = json.Unmarshal(resBody, &r)
	if err != nil {
		return nil, err
	}

	r = parseJson(r)
	return r, nil
}

func parseJson(r interface{}) interface{} {
	switch t := r.(type) {
	case []interface{}:
		var a []interface{}
		for _, item := range t {
			a = append(a, parseJson(item))
		}
		return a
	case map[string]interface{}:
		return t
	default:
		return t
	}
}

func sign(message, secret *string) string {
	mac := hmac.New(sha256.New, []byte(*secret))
	mac.Write([]byte(*message))
	signature := fmt.Sprintf("%x", mac.Sum(nil))
	return signature
}
