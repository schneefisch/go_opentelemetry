package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type Parent struct {
	TraceID    string `json:"TraceID"`
	SpanID     string `json:"SpanID"`
	TraceFlags string `json:"TraceFlags"`
	TraceState string `json:"TraceState"`
}

type SpanContext struct {
	TraceID    string `json:"TraceID"`
	SpanID     string `json:"SpanID"`
	TraceFlags string `json:"TraceFlags"`
	TraceState string `json:"TraceState"`
}

type TraceType struct {
	Name        string      `json:"Name"`
	Parent      Parent      `json:"Parent"`
	SpanContext SpanContext `json:"SpanContext"`
}

func Test_testHandler(t *testing.T) {

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/test", nil)

	myRequestHandler(recorder, request)
	closeTrace()

	if recorder.Code != http.StatusOK {
		t.Errorf("Unexpected Error-Response")
	}

	// check if traceResult-file exists
	traceFile, err := os.Stat("traces.txt")
	if errors.Is(err, os.ErrNotExist) {
		t.Errorf("Traces file not found, but expected one: %v", err)
	}
	// check if tracefile is empty
	if traceFile.Size() < 1 {
		t.Errorf("File size is too small! size = %d", traceFile.Size())
	}

	// read file
	file, err := os.ReadFile("traces.txt")
	if err != nil {
		t.Errorf("Could not read file: %v", err)
	}
	var traceResult TraceType
	err = json.Unmarshal(file, &traceResult)
	if err != nil {
		t.Errorf("Could not unmarshall traceResult: %v", err)
	}

	// expecting Parent.TraceID to be set
	if traceResult.Parent.TraceID == "" {
		t.Errorf("Missing Parent.TraceID value")
	}
	// expecting SpanContext.TraceID to be set
	if traceResult.SpanContext.TraceID == "" {
		t.Errorf("Missing SpanContext.TraceID value")
	}
}

func Test_testWithTraceParentInHeader(t *testing.T) {

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/test", nil)

	request.Header.Set("traceparent", "")

	myRequestHandler(recorder, request)
	closeTrace()

	if recorder.Code != http.StatusOK {
		t.Errorf("Unexpected Error-Response")
	}

	// check if traceResult-file exists
	traceFile, err := os.Stat("traces.txt")
	if errors.Is(err, os.ErrNotExist) {
		t.Errorf("Traces file not found, but expected one: %v", err)
	}
	// check if tracefile is empty
	if traceFile.Size() < 1 {
		t.Errorf("File size is too small! size = %d", traceFile.Size())
	}

	// read file
	file, err := os.ReadFile("traces.txt")
	if err != nil {
		t.Errorf("Could not read file: %v", err)
	}
	var traceResult TraceType
	err = json.Unmarshal(file, &traceResult)
	if err != nil {
		t.Errorf("Could not unmarshall traceResult: %v", err)
	}

	// expecting Parent.TraceID to be set
	if traceResult.Parent.TraceID == "" {
		t.Errorf("Missing Parent.TraceID value")
	}
	// expecting SpanContext.TraceID to be set
	if traceResult.SpanContext.TraceID == "" {
		t.Errorf("Missing SpanContext.TraceID value")
	}
}
