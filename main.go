package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		method := r.Method
		path := r.URL.Path
		body, err := io.ReadAll(r.Body)
		header := r.Header

		if err != nil {
			if err.Error() != "EOF" {
				log.Fatal(err)
			}
		}

		var bodyRequest, statusRequest = conectionApi(method, path, body, header)

		w.WriteHeader(statusRequest)
		fmt.Fprintf(w, bodyRequest)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func conectionApi(method string, path string, body []byte, header http.Header) (string, int) {
	client := http.Client{}

	var bodyContent io.Reader

	if len(body) == 0 {
		bodyContent = nil
	} else {
		bodyContent = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, "http://localhost:4567"+path, bodyContent)

	for key, value := range header {
		req.Header.Set(key, strings.Join(value, ""))
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return string(body), resp.StatusCode
}
