package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"starter/app/helpers"
	"starter/app/models"
	"starter/app/utils"
)

// Register a new user
func (app *App) Register(w http.ResponseWriter, r *http.Request) {
	var response = map[string]interface{}{"status": "success", "message": "User registered successfully"}

	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
	}

	usr, _ := user.GetUser(app.DB)
	if usr != nil {
		response["status"] = "failed"
		response["message"] = "User already registered"
		helpers.JSON(w, http.StatusBadRequest, response)
		return
	}

	// Strip the next of white space
	user.Prepare()

	// Register Validation
	err = user.Validate("")
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	userCreated, err := user.SaveUser(app.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
	}
	response["user"] = userCreated
	helpers.JSON(w, http.StatusCreated, response)
}

// Login
func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	var response = map[string]interface{}{"status": "success", "message": "Login success"}

	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Remove whitespace
	user.Prepare()

	// Validate
	err = user.Validate("login")
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	userData, err := user.GetUser(app.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check user data
	if userData == nil {
		response["status"] = "error"
		response["message"] = "Account not recognized, please register"
		helpers.JSON(w, http.StatusUnauthorized, response)
		return
	}

	// Check the password
	err = models.CheckPasswordHash(user.Password, userData.Password)
	if err != nil {
		response["status"] = "error"
		response["message"] = "Invalid credentials"
		helpers.JSON(w, http.StatusUnprocessableEntity, response)
		return
	}

	token, err := utils.EncodeAuthToken(userData.ID)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	response["token"] = token
	helpers.JSON(w, http.StatusOK, response)
}

// GetAllUsers getting all users
func (app *App) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"status": "success", "message": "Users List"}

	user := &models.User{}

	users, err :=  user.GetUsers(app.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response["data"] = users
	helpers.JSON(w, http.StatusOK, response)
}