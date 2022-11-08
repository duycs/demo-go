package controllers

import (
	"github.com/codegangsta/negroni"
	"github.com/duycs/demo-go/demo/application/services"
	"github.com/gorilla/mux"
)

func RegisterUserHandlers(r *mux.Router, n negroni.Negroni, service services.UserUseCase) {
	r.Handle("/v1/users", n.With(
		negroni.Wrap(ListUsers(service)),
	)).Methods("GET", "OPTIONS").Name("listUsers")

	r.Handle("/v1/users/{id}", n.With(
		negroni.Wrap(GetUser(service)),
	)).Methods("GET", "OPTIONS").Name("getUser")

	r.Handle("/v1/users", n.With(
		negroni.Wrap(CreateUser(service)),
	)).Methods("POST", "OPTIONS").Name("createUser")

	r.Handle("/v1/users/{id}", n.With(
		negroni.Wrap(DeleteUser(service)),
	)).Methods("DELETE", "OPTIONS").Name("deleteUser")
}

func RegisterTaskHandlers(r *mux.Router, n negroni.Negroni, service services.TaskUseCase) {
	r.Handle("/v1/tasks", n.With(
		negroni.Wrap(ListTasks(service)),
	)).Methods("GET", "OPTIONS").Name("listTasks")

	r.Handle("/v1/tasks", n.With(
		negroni.Wrap(CreateTask(service)),
	)).Methods("POST", "OPTIONS").Name("createTask")

	r.Handle("/v1/tasks/{id}", n.With(
		negroni.Wrap(GetTask(service)),
	)).Methods("GET", "OPTIONS").Name("getTask")

	r.Handle("/v1/tasks/{id}", n.With(
		negroni.Wrap(DeleteTask(service)),
	)).Methods("DELETE", "OPTIONS").Name("deleteTask")
}
