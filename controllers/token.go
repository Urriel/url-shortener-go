package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urriel/url-shortener-go/models"
	"github.com/urriel/url-shortener-go/utils"
	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// NewTokenController create a controller and assign the routes
func NewTokenController(mux *web.Mux, database *mgo.Database) {
	fmt.Println("Registering TokenController")

	controller := tokenController{
		DB: database,
	}

	mux.Get("/tokens", controller.Index)
	mux.Post("/tokens/create", controller.Create)
	mux.Get("/t/:token", controller.Redirect)
}

// === Token Controller ===

type tokenController struct {
	DB *mgo.Database
}

// === Create Token ===

// createRequestPayload should handle the request body
type createRequestPayload struct {
	URL string `json:"url"`
}

// createResponsePayload used to marshal the response
type createResponsePayload struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

func (controller *tokenController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	buff := utils.ReadBody(r)
	body := new(createRequestPayload)
	err := json.Unmarshal(buff, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSONError(w, "JSON payload invalid")
		return
	}

	generatedToken, err := utils.GenerateToken()
	if err != nil {
		_ = fmt.Errorf("Cannot generate a new token : %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJSONError(w, "Cannot generate a new token")
		return
	}

	err = controller.DB.C("tokens").Insert(&models.Token{
		URL:    body.URL,
		Token:  generatedToken,
		Visits: 0,
	})
	if err != nil {
		fmt.Errorf("Cannot save the new token in the database : %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteJSONError(w, "Cannot save the new token")
		return
	}

	payload := &createResponsePayload{
		Token: generatedToken,
		URL:   body.URL,
	}
	utils.WriteJSONSuccess(w, payload)
}

// === Index Token ===

type indexResponsePayload struct {
	Data []models.Token
}

func (controller *tokenController) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	tokens := []models.Token{}

	err := controller.DB.C("tokens").Find(bson.M{}).All(&tokens)
	if err != nil {
		fmt.Errorf("Cannot fetch the list of tokens : %s", err.Error())

		answerPayload := &utils.HTTPError{
			Msg:  "Cannot fetch the list of tokens",
			Code: "TOKEN02",
		}

		answer, _ := json.Marshal(answerPayload)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(answer)
	}

	answerPayload := &indexResponsePayload{
		Data: tokens,
	}
	answer, _ := json.Marshal(answerPayload)
	w.Write(answer)
}

// === Redirect Token ===

func (controller *tokenController) Redirect(c web.C, w http.ResponseWriter, r *http.Request) {
	tokenModel := new(models.Token)
	tokenString := c.URLParams["token"]

	err := controller.DB.C("tokens").Find(bson.M{"token": tokenString}).One(tokenModel)
	if err != nil {
		fmt.Errorf("Cannot find the right token : %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry we cannot forward you to the right url, try again later"))
		return
	}

	if len(tokenModel.URL) == 0 {
		fmt.Println("Token not found")

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Redirection not found"))
		return
	}

	// Increment the number of visits
	err = controller.DB.C("tokens").Update(bson.M{"token": tokenString}, bson.M{"$inc": bson.M{"visits": 1}})
	if err != nil {
		fmt.Errorf("Cannot increment the visit counter : %s", err.Error())
	}

	http.Redirect(w, r, tokenModel.URL, http.StatusPermanentRedirect)
}
