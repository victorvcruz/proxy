package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			if err.Error() != "EOF" {
				log.Fatal(err)
			}
		}

		method := r.Method
		path := r.URL.Path
		header := r.Header

		response := requestToAPI(method, path, body, header)

		for key, value := range response.headers {
			w.Header().Set(key, strings.Join(value, ""))
		}

		w.WriteHeader(response.status)
		fmt.Fprintf(w, response.body)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}

type RequisitionInformation struct {
	Path        string
	Count       int64
	Duration    int64
	AverageTime int64
}

type ResponseAPI struct {
	body    string
	status  int
	headers map[string][]string
}

var apiRequisitonMap = make(map[string]RequisitionInformation)

func requestToAPI(method string, path string, bodyRequest []byte, header http.Header) ResponseAPI {
	client := http.Client{}

	var bodyContent io.Reader

	if len(bodyRequest) == 0 {
		bodyContent = nil
	} else {
		bodyContent = bytes.NewBuffer(bodyRequest)
	}

	req, err := http.NewRequest(method, "http://localhost:4567"+path, bodyContent)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range header {
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

	pathID := req.Method + "-" + req.URL.Path
	elapsed := time.Since(start)

	if val, ok := apiRequisitonMap[pathID]; ok {
		val.Count = val.Count + 1
		val.Duration += elapsed.Milliseconds()
		val.AverageTime = val.Duration / val.Count
		apiRequisitonMap[pathID] = val
	} else {
		var m RequisitionInformation
		m.Path = req.URL.Path
		m.Count = 1
		m.Duration = elapsed.Milliseconds()
		m.AverageTime = m.Duration / m.Count
		apiRequisitonMap[pathID] = m
	}

	resp.Header.Set("Count", strconv.FormatInt(apiRequisitonMap[pathID].Count, 10))
	resp.Header.Set("Duration", strconv.FormatInt(apiRequisitonMap[pathID].Duration, 10)+"ms")
	resp.Header.Set("AverageTime", strconv.FormatInt(apiRequisitonMap[pathID].AverageTime, 10)+"ms")

	return ResponseAPI{
		body:    string(bodyResponse),
		status:  resp.StatusCode,
		headers: resp.Header,
	}
}
