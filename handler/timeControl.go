package handler

import (
	"net/http"
	"strconv"
	"time"
)

type RequisitionInformation struct {
	Path        string
	Count       int64
	Duration    int64
	AverageTime int64
}

var apiRequisitonMap = make(map[string]RequisitionInformation)

func HandlerTimeAndInsertHeaders(resp *http.Response, req *http.Request, start time.Time) {
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
}
