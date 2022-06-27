package handler

import (
	"encoding/json"
	"net/http"
	"proxy_project/api/response"
	"proxy_project/cache"
)

func InsertCache(client cache.CacheClient, req *http.Request, query string, responseApi *response.ResponseAPI) error {
	reqID := req.Method + "-" + req.URL.Path + query + "-" + req.Header.Get("Token")

	respJson, err := json.Marshal(responseApi)
	if err != nil {
		return err
	}

	err = client.InsertInDatabase(reqID, string(respJson))
	if err != nil {
		return err
	}

	return nil
}

func InsertCacheArray(client cache.CacheClient, req *http.Request, query string, responseApi *response.ResponseAPIArray) error {
	reqID := req.Method + "-" + req.URL.Path + query + "-" + req.Header.Get("Token")

	respJson, err := json.Marshal(responseApi)
	if err != nil {
		return err
	}

	err = client.InsertInDatabase(reqID, string(respJson))
	if err != nil {
		return err
	}

	return nil
}

func FindInCache(client cache.CacheClient, req *http.Request, query string) (*response.ResponseAPI, error) {
	reqID := req.Method + "-" + req.URL.Path + query + "-" + req.Header.Get("Token")

	val, err := client.FindInDatabase(reqID)
	if err != nil {
		switch e := err.(type) {
		case *cache.CacheNotFoundError:
			return nil, e
		}
	}

	var responseCache response.ResponseAPI

	err = json.Unmarshal([]byte(val), &responseCache)
	if err != nil {
		return nil, err
	}

	return &responseCache, nil
}

func FindInCacheArray(client cache.CacheClient, req *http.Request, query string) (*response.ResponseAPIArray, error) {
	reqID := req.Method + "-" + req.URL.Path + query + "-" + req.Header.Get("Token")

	val, err := client.FindInDatabase(reqID)
	if err != nil {
		switch e := err.(type) {
		case *cache.CacheNotFoundError:
			return nil, e
		}
	}

	var responseCache response.ResponseAPIArray

	err = json.Unmarshal([]byte(val), &responseCache)
	if err != nil {
		return nil, err
	}

	return &responseCache, nil
}
