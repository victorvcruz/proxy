package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ResponseAPI struct {
	Body    map[string]interface{}
	Status  int
	Headers map[string][]string
}

func (r *ResponseAPI) ResponseWriter(w http.ResponseWriter) error {
	for key, value := range r.Headers {
		w.Header().Set(key, strings.Join(value, ""))
	}

	bodyJson, err := json.Marshal(r.Body)
	if err != nil {
		return err
	}

	w.WriteHeader(r.Status)
	fmt.Fprintf(w, string(bodyJson))
	return nil
}

type ResponseAPIArray struct {
	Body    []interface{}
	Status  int
	Headers map[string][]string
}

func (r *ResponseAPIArray) ResponseWriter(w http.ResponseWriter) error {
	for key, value := range r.Headers {
		w.Header().Set(key, strings.Join(value, ""))
	}

	bodyJson, err := json.Marshal(r.Body)
	if err != nil {
		return err
	}

	w.WriteHeader(r.Status)
	fmt.Fprintf(w, string(bodyJson))
	return nil
}
