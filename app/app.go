package app

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"io"
	"net/http"
	"os"
)

const name = "otel"

func newConsoleExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// human readable output
		stdouttrace.WithPrettyPrint(),
		// to not print timestamp for demo
		stdouttrace.WithoutTimestamps(),
	)
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
		),
	)
	return r
}

var traceProvider *trace.TracerProvider

func initTrace() {
	// write telemetry-data to file
	file, err := os.Create("traces.txt")
	if err != nil {
		log.Err(err)
	}

	exp, err := newConsoleExporter(file)
	if err != nil {
		log.Err(err)
	}

	traceProvider = trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	otel.SetTracerProvider(traceProvider)
}

func closeTrace() {
	traceProvider.ForceFlush(context.Background())
	if err := traceProvider.Shutdown(context.Background()); err != nil {
		log.Err(err)
	}
}

func InitApp() *chi.Mux {

	// init endpoints and server
	router := chi.NewRouter()
	router.Get("/api/test", myRequestHandler)

	return router
}

func myRequestHandler(writer http.ResponseWriter, request *http.Request) {
	// adding tracer
	_, span := otel.Tracer(name).Start(request.Context(), "Run")
	defer span.End()

	// do something complicated
	in := uint(10)
	number, err := Fibonacci(in)
	if err != nil {
		// log something
		log.Info().Msgf("Fibonacci(%d): %v\n", in, err)
	} else {
		log.Info().Msgf("Fibonacci(%d): %d\n", in, number)
	}

	// write responses
	numberString, _ := json.Marshal(number)
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(numberString)
}

func Fibonacci(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}
	var n2, n1 uint64 = 0, 1
	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}
	return n2 + n1, nil
}
