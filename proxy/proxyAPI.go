package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"proxy_project/api"
	"proxy_project/cache"
	"proxy_project/handler"
	"strings"
)

type ProxyAPI struct {
	cache.CacheClient
	api.RequestClient
}

func (p *ProxyAPI) Run() error {

	var mutex handler.QueueMutex

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			if err.Error() != "EOF" {
				log.Fatal(err)
			}
		}

		queryParms := transformMapInQueryParams(r.URL.Query())
		log.Println(r.Method)

		unlock := mutex.Lock(r, body, queryParms)
		defer unlock()

		if r.Method == "GET" {
			responseCache, err := handler.FindInCache(p.CacheClient, r, queryParms)

			switch e := err.(type) {
			case *cache.CacheNotFoundError:
				log.Println(e)

			case *json.UnmarshalTypeError:
				responseCache, err := handler.FindInCacheArray(p.CacheClient, r, queryParms)
				log.Println("Array cache")

				switch e := err.(type) {
				case *cache.CacheNotFoundError:
					log.Println(e)

				default:
					if err := responseCache.ResponseWriter(w); err != nil {
						log.Fatal(err)
					}

					return
				}
			default:
				if err := responseCache.ResponseWriter(w); err != nil {
					log.Fatal(err)
				}
				log.Println("Cache")

				return

			}
		}

		log.Println("Requisition")

		response, responseArray := p.RequestClient.RequestAPI(p.CacheClient, r, body, queryParms)

		if response != nil {
			if err := (*response).ResponseWriter(w); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := (*responseArray).ResponseWriter(w); err != nil {
				log.Fatal(err)
			}
		}

		return
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}

	return nil

}

func transformMapInQueryParams(query map[string][]string) string {
	if len(query) == 0 {
		return ""
	}

	queryContent := new(bytes.Buffer)
	fmt.Fprintf(queryContent, "?")
	for key, value := range query {
		fmt.Fprintf(queryContent, "%s=%s&", key, strings.Join(value, ""))
	}

	return queryContent.String()[:len(queryContent.String())-1]
}
