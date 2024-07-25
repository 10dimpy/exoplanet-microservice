package handlers

import (
	"encoding/json"
	"exoplanet-microservice/models"
	"exoplanet-microservice/services"
	"exoplanet-microservice/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	Service services.ExoplanetService
}

func (h *Handler) AddExoplanet(w http.ResponseWriter, r *http.Request) {
	var exoplanet models.Exoplanet
	if err := json.NewDecoder(r.Body).Decode(&exoplanet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.ValidateExoplanet(exoplanet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exoplanet.ID = h.Service.AddExoplanet(exoplanet)
	json.NewEncoder(w).Encode(exoplanet)
}

func (h *Handler) ListExoplanets(w http.ResponseWriter, r *http.Request) {
	exoplanets := h.Service.ListExoplanets()
	json.NewEncoder(w).Encode(exoplanets)
}

func (h *Handler) GetExoplanetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	exoplanet, err := h.Service.GetExoplanetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(exoplanet)
}

func (h *Handler) UpdateExoplanet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	var exoplanet models.Exoplanet
	if err := json.NewDecoder(r.Body).Decode(&exoplanet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	exoplanet.ID = id

	if err := utils.ValidateExoplanet(exoplanet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedExoplanet, err := h.Service.UpdateExoplanet(exoplanet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(updatedExoplanet)
}

func (h *Handler) DeleteExoplanet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if err := h.Service.DeleteExoplanet(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) EstimateFuel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	crewCapacity, _ := strconv.Atoi(r.URL.Query().Get("crew"))

	fuel, err := h.Service.EstimateFuel(id, crewCapacity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]float64{"fuel": fuel})
}
