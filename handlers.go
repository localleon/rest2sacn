package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func sacnHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	log.Println(vars)
	w.WriteHeader(http.StatusOK)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {

}
