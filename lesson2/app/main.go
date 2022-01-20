package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type app struct {
	logger *zap.Logger
}

func (a *app) processHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.With(zap.Any("method", r.Method)).Debug("process handler called")
	if r.Method != http.MethodGet {
		a.logger.With(zap.Any("method_not_allowed", r.Method)).Debug("process handler called")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("ID")

	if len(id) > 0 {
		a.writeResponseByPersonId(w, id)
		return
	}

	a.writeResponsePersons(w)
}

func (a *app) Serve() error {
	mux := http.NewServeMux()
	mux.Handle("/significant_persons", http.HandlerFunc(a.processHandler))

	return http.ListenAndServe("0.0.0.0:8080", mux)
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	a := app{
		logger: logger,
	}

	a.logger.Debug("starting server ...")
	if err := a.Serve(); err != nil {
		a.logger.With(zap.Error(err)).Fatal("couldnt start server")
	}
}

func (a *app) writeResponsePersons(w http.ResponseWriter) {
	a.logger.With(zap.Bool("all_persons", true)).Debug("response all persons")

	persons := getPersons()

	w.WriteHeader(http.StatusOK)
	err := persons.ToJson(w)
	if err != nil {
		http.Error(w, "Unable to process request", http.StatusInternalServerError)
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func (a *app) writeResponseByPersonId(w http.ResponseWriter, id string) {
	a.logger.With(zap.Any("id", id)).Debug("response by person id")

	person := getPersonById(id)

	w.WriteHeader(http.StatusOK)
	err := person.ToJson(w)
	if err != nil {
		a.logger.Error("writeResponseByPersonId couldnt convert to json", zap.Error(err))
		http.Error(w, "Unable to process request", http.StatusInternalServerError)
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}
