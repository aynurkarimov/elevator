package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/aynurkarimov/elevator/internal/action"
	"github.com/aynurkarimov/elevator/internal/elevator"
	"github.com/google/uuid"
)

type RequestElevatorDto struct {
	ElevatorId     uuid.UUID `json:"elevatorId"`
	RequestedFloor uint8     `json:"requestedFloor"`
}

func Handler() error {
	elevators := elevator.Elevators
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("ui", "templates", "index.html")
		t, err := template.ParseFiles(path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		t.ExecuteTemplate(w, "index", elevators)
	})

	mux.HandleFunc("/request-elevator", func(w http.ResponseWriter, r *http.Request) {
		var dto RequestElevatorDto

		err := json.NewDecoder(r.Body).Decode(&dto)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		teardown, err := action.PerformAction(&action.Action{
			ElevatorId: dto.ElevatorId,
			Command:    action.CommandReq,
			Floor:      dto.RequestedFloor,
		})

		data := map[string]int{"teardown": int(teardown)}
		jsonData, err := json.Marshal(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	mux.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("ui", "templates", "elevators.html")
		t, err := template.ParseFiles(path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		t.ExecuteTemplate(w, "elevators", elevators)
	})

	err := http.ListenAndServe(":3000", mux)

	if err != nil {
		return fmt.Errorf("Couldn't listen and server :3000\n")
	}

	return nil
}
