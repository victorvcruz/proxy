package requestAPI

import (
	"net/http"
	"proxy_project/cache"
	"proxy_project/proxyAPI/response"
)

type RequestClient interface {
	RequestToAPI(cacheClient cache.CacheClient, r *http.Request, bodyRequest []byte, query string) (**response.ResponseAPI, **response.ResponseAPIArray)
}
