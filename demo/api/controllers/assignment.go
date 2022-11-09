package controllers

import (
	"fmt"
	"net/http"

	"github.com/duycs/demo-go/demo/application/services"
	"github.com/duycs/demo-go/demo/infrastructure/helpers"
	"github.com/gorilla/mux"
)

func AssignTask(taskService services.TaskUseCase, userService services.UserUseCase, assignmentService services.AssignmentUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error assign task"

		vars := mux.Vars(r)
		taskID, err := helpers.StringToID(vars["task_id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		task, err := taskService.GetTask(taskID)
		if err != nil && err != helpers.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if task == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}

		userID, err := helpers.StringToID(vars["user_id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		user, err := userService.GetUser(userID)
		if err != nil && err != helpers.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if user == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}

		err = assignmentService.Assign(user, task)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

func Checkout(taskService services.TaskUseCase, assignmentService services.AssignmentUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error checkout task"

		vars := mux.Vars(r)
		taskID, err := helpers.StringToID(vars["task_id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		task, err := taskService.GetTask(taskID)
		if err != nil && err != helpers.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if task == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}

		err = assignmentService.Checkout(task)
		if err != nil && err != helpers.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}
