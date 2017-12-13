package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpQuery map[string]string

type HttpResponse struct {
	Body  []byte
	Error error
}

type ResponseChan chan HttpResponse

func (ch *ResponseChan) Read() ([]byte, error) {
	raw := <-*ch
	return raw.Body, raw.Error
}

func (ch *ResponseChan) AsJSON(ref interface{}) error {
	body, err := ch.Read()

	if err != nil {
		return err
	} else {
		return json.Unmarshal(body, ref)
	}
}

func HttpGet(u string, query HttpQuery) *ResponseChan {
	ch := make(ResponseChan, 1)

	go func() {
		values := url.Values{}

		for k, v := range query {
			values.Add(k, v)
		}

		result := HttpResponse{}
		response, err := http.Get(u + "?" + values.Encode())

		if err != nil {
			result.Error = err
		} else {
			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)

			if err != nil {
				result.Error = err
			} else {
				result.Body = body
			}
		}

		ch <- result
	}()

	return &ch
}
