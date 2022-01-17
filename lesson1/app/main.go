package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	Namespace = "ocmetricsexample"

	LabelMethod = "method"
	LabelStatus = "status"
)

type app struct {
	latencyHistogram,
	personsHistogram *prometheus.HistogramVec
	personsCounter prometheus.Counter
}

func (a *app) processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	startTime := time.Now()

	id := r.URL.Query().Get("ID")

	defer func() {
		a.personsCounter.Inc()

		a.latencyHistogram.With(prometheus.Labels{LabelMethod: r.Method}).Observe(sinceInMilliseconds(startTime))

		id, err := strconv.ParseFloat(id, 64)
		if err == nil {
			a.personsHistogram.With(prometheus.Labels{LabelStatus: "OK"}).Observe((id))
		}
	}()

	if len(id) > 0 {
		writeResponseByPersonId(w, id)
		return
	}
	writeResponsePersons(w)

}

func (a *app) Init() error {
	a.personsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "persons_counter",
		Help:      "The number of times persons are requested",
	})

	prometheus.MustRegister(a.personsCounter)

	a.latencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: Namespace,
		Name:      "persons_latency",
		Help:      "The distribution of the latencies",
		Buckets:   []float64{0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000},
	}, []string{LabelMethod})

	prometheus.MustRegister(a.latencyHistogram)

	a.personsHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: Namespace,
		Name:      "persons_by_ids",
		Help:      "Groups persons requests by ids",
		Buckets:   []float64{1, 2, 3, 4, 5},
	}, []string{LabelStatus})

	prometheus.MustRegister(a.personsHistogram)

	return nil
}

func (a *app) Serve() error {
	mux := http.NewServeMux()
	mux.Handle("/significant_persons", http.HandlerFunc(a.processHandler))
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe("0.0.0.0:9000", mux)
}

func main() {
	a := app{}

	if err := a.Init(); err != nil {
		log.Fatal(err)
	}

	if err := a.Serve(); err != nil {
		log.Fatal(err)
	}

}

func sinceInMilliseconds(startTime time.Time) float64 {
	return float64(time.Since(startTime).Nanoseconds()) / 1e6
}

func writeResponsePersons(w http.ResponseWriter) {
	persons := getPersons()

	w.WriteHeader(http.StatusOK)
	err := persons.ToJson(w)
	if err != nil {
		http.Error(w, "Unable to process request", http.StatusInternalServerError)
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func writeResponseByPersonId(w http.ResponseWriter, id string) {
	person := getPersonById(id)

	w.WriteHeader(http.StatusOK)
	err := person.ToJson(w)
	if err != nil {
		http.Error(w, "Unable to process request", http.StatusInternalServerError)
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}
