package app

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_testHandler(t *testing.T) {

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	//request.Header.Add()

	//rootTrace
	//request.Header.Set("traceparen", rootTrace)

	rctx := chi.NewRouteContext()
	// adding a uri-parameter
	//rctx.URLParams.Add("key", "value")

	// setting chi route-context
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	testHandler(recorder, request)
	closeTrace()

	if recorder.Code != http.StatusOK {
		t.Errorf("Unexpected Error-Response")
	}

	// check if trace-file exists
	trace, err := os.Stat("traces.txt")
	if errors.Is(err, os.ErrNotExist) {
		t.Errorf("Traces file not found, but expected one")
	}
	// check if tracefile is empty
	if trace.Size() < 1 {
		t.Errorf("File size is too small! size = %d", trace.Size())
	}
}
