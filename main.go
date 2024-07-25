package main

import (
	"exoplanet-microservice/handlers"
	"exoplanet-microservice/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	service := services.NewService()
	handler := &handlers.Handler{Service: service}

	r := mux.NewRouter()
	r.HandleFunc("/exoplanets", handler.AddExoplanet).Methods("POST")
	r.HandleFunc("/exoplanets", handler.ListExoplanets).Methods("GET")
	r.HandleFunc("/exoplanets/{id}", handler.GetExoplanetByID).Methods("GET")
	r.HandleFunc("/exoplanets/{id}", handler.UpdateExoplanet).Methods("PUT")
	r.HandleFunc("/exoplanets/{id}", handler.DeleteExoplanet).Methods("DELETE")
	r.HandleFunc("/exoplanets/{id}/fuel", handler.EstimateFuel).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
