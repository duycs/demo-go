package controllers

import (
	"net/http"

	"github.com/duycs/demo-go/demo"
	"github.com/duycs/demo-go/demo/infratructure/helpers"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	helpers.JSON(w, http.StatusOK, "Welcome To This Awesome API")
}
