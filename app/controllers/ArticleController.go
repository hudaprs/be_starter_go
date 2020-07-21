package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"starter/app/helpers"
	"starter/app/models"
)

func (app *App) CreateArticle(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"status": "success", "message": "Article created"}
	article := &models.Article{}
	user := r.Context().Value("userID").(float64)


	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &article)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	article.Prepare()

	err = article.Validate()
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	article.UserID = uint(user)

	articleData, err := article.Save(app.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
	}

	response["data"] = articleData
	helpers.JSON(w,	 http.StatusCreated, response)
	return
}

func (app *App) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"Status": "success", "Message": "Article Deleted"}
	article := &models.Article{}
	id := mux.Vars(r)["id"]

	article.DeleteArticle(id, app.DB)
	helpers.JSON(w, http.StatusOK, response)
}
