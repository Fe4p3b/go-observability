package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

type app struct {
	logger *zap.Logger
	tracer opentracing.Tracer
}

type zapWrapper struct {
	logger *zap.Logger
}

// Error logs a message at error priority
func (w *zapWrapper) Error(msg string) {
	w.logger.Error(msg)
}

// Infof logs a message at info priority
func (w *zapWrapper) Infof(msg string, args ...interface{}) {
	w.logger.Sugar().Infof(msg, args...)
}

func initJaeger(service string, logger *zap.Logger) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(&zapWrapper{logger: logger}))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return tracer, closer
}

func (a *app) processHandler(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(r.Context(), a.tracer, "processHandler")
	defer span.Finish()

	span.LogFields(
		log.String("method", r.Method),
	)

	a.logger.With(zap.Any("method", r.Method)).Debug("process handler called")
	if r.Method != http.MethodGet {
		a.logger.With(zap.Any("method_not_allowed", r.Method)).Debug("process handler called")
		span.LogFields(
			log.Error(fmt.Errorf("method not allowed")),
		)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	id := r.URL.Query().Get("ID")

	if len(id) > 0 {
		a.writeResponseByPersonId(ctx, w, id)
		return
	}

	a.writeResponsePersons(ctx, w)
}

func (a *app) Serve() error {
	http.Handle("/significant_persons", http.HandlerFunc(a.processHandler))

	return http.ListenAndServe("0.0.0.0:9000", nethttp.Middleware(a.tracer, http.DefaultServeMux))
}

func (a *app) writeResponsePersons(ctx context.Context, w http.ResponseWriter) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, a.tracer, "writeResponsePersons")
	defer span.Finish()

	a.logger.With(zap.Bool("all_persons", true)).Debug("response all persons")

	persons := getPersons()

	r := rand.Intn(10)
	a.logger.Debug("random generated integer ", zap.Int("rand", r))
	if r > 5 {
		span.LogFields(
			log.Error(fmt.Errorf("error on purpose")),
		)

		http.Error(w, "Unable to process request", http.StatusInternalServerError)
		return
	}

	span.LogFields(
		log.String("status", "StatusOK"),
	)

	w.WriteHeader(http.StatusOK)
	err := persons.ToJson(w)
	if err != nil {
		span.LogFields(
			log.Error(err),
		)

		http.Error(w, "Unable to process request", http.StatusInternalServerError)
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func (a *app) writeResponseByPersonId(ctx context.Context, w http.ResponseWriter, id string) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, a.tracer, "writeResponseByPersonId")
	defer span.Finish()

	a.logger.With(zap.Any("id", id)).Debug("response by person id")

	person := getPersonById(id)

	span.LogFields(
		log.String("status", "StatusOK"),
		log.String("arg0", id),
	)

	w.WriteHeader(http.StatusOK)
	err := person.ToJson(w)
	if err != nil {
		a.logger.Error("writeResponseByPersonId couldnt convert to json", zap.Error(err))

		span.LogFields(
			log.Error(err),
		)

		http.Error(w, "Unable to process request", http.StatusInternalServerError)
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}
