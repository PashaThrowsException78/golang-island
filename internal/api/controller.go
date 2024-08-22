package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"golang-island/internal/dto"
	"golang-island/internal/service"
	"golang.org/x/exp/slog"
	"net/http"
	"strconv"
)

type IslandController struct {
	service service.IslandService
	log     *slog.Logger
}

func NewController(log *slog.Logger) *IslandController {
	return &IslandController{service: service.NewService(log), log: log}
}

func (controller IslandController) CalculateIsland(w http.ResponseWriter, r *http.Request) {

	var request dto.CalculateIslandsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		controller.log.Warn("Error decoding request", "error", err)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	err = controller.service.PutTask(request)

	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		controller.log.Warn("Unable to create task id=", request.IslandId, " ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	}
}

func (controller IslandController) GetIslandResult(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	islandCount, err := controller.service.GetResult(id)

	response := map[string]int{"islandCount": islandCount}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		controller.log.Warn("Error encoding response", "error", err)
	}
}

func (controller IslandController) IsReady(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	ready, err := controller.service.IsReady(id)

	exists := err == nil

	response := map[string]bool{"ready": ready, "exists": exists}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		controller.log.Warn("Error encoding response", "error", err)
	}
}
