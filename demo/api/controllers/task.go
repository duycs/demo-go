package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/duycs/demo-go/demo/application/dto"
	"github.com/duycs/demo-go/demo/application/services"
	"github.com/duycs/demo-go/demo/domain/entity"
	"github.com/duycs/demo-go/demo/infrastructure/helpers"
	"github.com/gorilla/mux"
)

func ListTasks(service services.TaskUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading tasks"

		var data []*entity.Task
		var err error
		title := r.URL.Query().Get("title")

		switch {
		case title == "":
			data, err = service.ListTasks()
		default:
			data, err = service.SearchTasks(title)
		}

		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != helpers.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}

		var toJ []*dto.Task
		for _, d := range data {
			toJ = append(toJ, &dto.Task{
				ID:                 d.ID,
				Title:              d.Title,
				Description:        d.Description,
				EstimationInSecond: d.EstimationInSecond,
			})
		}

		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func CreateTask(service services.TaskUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding task"

		var input struct {
			Title              string `json:"title"`
			Description        string `json:"description"`
			EstimationInSecond int    `json:"estimation_in_second"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		id, err := service.CreateTask(input.Title, input.Description, input.EstimationInSecond)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		toJ := &dto.Task{
			ID:                 id,
			Title:              input.Title,
			Description:        input.Description,
			EstimationInSecond: input.EstimationInSecond,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func GetTask(service services.TaskUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading task"
		vars := mux.Vars(r)
		id, err := helpers.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		data, err := service.GetTask(id)
		if err != nil && err != helpers.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}

		toJ := &dto.Task{
			ID:                 data.ID,
			Title:              data.Title,
			Description:        data.Description,
			EstimationInSecond: data.EstimationInSecond,
		}

		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func DeleteTask(service services.TaskUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing task"

		vars := mux.Vars(r)
		id, err := helpers.StringToID(vars["id"])

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		err = service.DeleteTask(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}
