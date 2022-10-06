package controllers

import (
	"net/http"

	"github.com/duycs/demo-go/demo/infrastructure/helpers"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	helpers.JSON(w, http.StatusOK, "Welcome to Home page")
}
