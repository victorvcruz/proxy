package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"proxy_project/cache"
	"proxy_project/handler"
	"proxy_project/proxyAPI/response"
	"strings"
	"time"
)

func RequestToAPI(cacheClient cache.CacheClient, r *http.Request, bodyRequest []byte, query string) (*response.ResponseAPI, *response.ResponseAPIArray) {
	client := http.Client{}
	var bodyContent io.Reader

	if len(bodyRequest) == 0 {
		bodyContent = nil
	} else {
		bodyContent = bytes.NewBuffer(bodyRequest)
	}

	req, err := http.NewRequest(r.Method, "http://localhost:4567"+r.URL.Path+query, bodyContent)
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

	var mapBodyResponse map[string]interface{}
	var mapBodyResponseArray []interface{}

	if err = json.Unmarshal(bodyResponse, &mapBodyResponse); err != nil {

		if err = json.Unmarshal(bodyResponse, &mapBodyResponseArray); err != nil {
			log.Fatal(err)
		}
		response := response.ResponseAPIArray{
			Body:    mapBodyResponseArray,
			Status:  resp.StatusCode,
			Headers: resp.Header,
		}

		if req.Method == "GET" && resp.StatusCode == 200 {
			if err := handler.HandlerInsertCacheArray(cacheClient, req, query, response); err != nil {
				log.Fatal(err)
			}
		}
		return nil, &response
	}

	response := response.ResponseAPI{
		Body:    mapBodyResponse,
		Status:  resp.StatusCode,
		Headers: resp.Header,
	}

	if req.Method == "GET" && resp.StatusCode == 200 {
		if err := handler.HandlerInsertCache(cacheClient, req, query, response); err != nil {
			log.Fatal(err)
		}
	}

	return &response, nil
}
