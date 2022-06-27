package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"proxy_project/api/response"
	"proxy_project/cache"
	"proxy_project/handler"
	"strings"
	"time"
)

type RequestClient struct {
	Host string
	Port string
}

func (request *RequestClient) RequestAPI(cacheClient cache.CacheClient, r *http.Request, bodyRequest []byte, query string) (**response.ResponseAPI, **response.ResponseAPIArray) {
	client := http.Client{}
	var bodyContent io.Reader

	if len(bodyRequest) == 0 {
		bodyContent = nil
	} else {
		bodyContent = bytes.NewBuffer(bodyRequest)
	}

	req, err := http.NewRequest(r.Method, request.Host+":"+request.Port+r.URL.Path+query, bodyContent)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range r.Header {
		req.Header.Set(key, strings.Join(value, ""))
	}

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bodyResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	handler.HandlerTimeAndInsertHeaders(resp, req, start)

	responseForProxy, err := CreateResponseForProxy(resp, bodyResponse)
	if err != nil {
		responseForProxy, err := CreateResponseArrayForProxy(resp, bodyResponse)
		if err != nil {
			log.Fatal(err)
		}

		if req.Method == "GET" && resp.StatusCode >= 200 && resp.StatusCode <= 299 {

			go func() {
				if err := handler.InsertCacheArray(cacheClient, req, query, responseForProxy); err != nil {
					log.Fatal(err)
				}
			}()

		}
		return nil, &responseForProxy
	}

	if req.Method == "GET" && resp.StatusCode >= 200 && resp.StatusCode <= 299 {

		go func() {
			if err := handler.InsertCache(cacheClient, req, query, responseForProxy); err != nil {
				log.Fatal(err)
			}
		}()

	}

	return &responseForProxy, nil
}

func CreateResponseForProxy(resp *http.Response, bodyResponse []byte) (*response.ResponseAPI, error) {

	var mapBodyResponse map[string]interface{}

	if err := json.Unmarshal(bodyResponse, &mapBodyResponse); err != nil {
		return nil, err
	}

	response := response.ResponseAPI{
		Body:    mapBodyResponse,
		Status:  resp.StatusCode,
		Headers: resp.Header,
	}

	return &response, nil

}

func CreateResponseArrayForProxy(resp *http.Response, bodyResponse []byte) (*response.ResponseAPIArray, error) {

	var mapBodyResponseArray []interface{}

	if err := json.Unmarshal(bodyResponse, &mapBodyResponseArray); err != nil {
		return nil, err
	}
	response := response.ResponseAPIArray{
		Body:    mapBodyResponseArray,
		Status:  resp.StatusCode,
		Headers: resp.Header,
	}

	return &response, nil

}
