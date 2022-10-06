package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/duycs/demo-go/API/demo/infrastructure/auth"
	"github.com/duycs/demo-go/API/demo/infrastructure/helpers"
	"github.com/duycs/demo-go/API/demo/models"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {

		formattedError := helpers.FormatError(err.Error())

		helpers.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	helpers.JSON(w, http.StatusCreated, userCreated)
}

func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	helpers.JSON(w, http.StatusOK, users)
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}
	helpers.JSON(w, http.StatusOK, userGotten)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		helpers.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		helpers.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := helpers.FormatError(err.Error())
		helpers.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	helpers.JSON(w, http.StatusOK, updatedUser)
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		helpers.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		helpers.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		helpers.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	helpers.JSON(w, http.StatusNoContent, "")
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := helpers.FormatError(err.Error())
		helpers.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	helpers.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
