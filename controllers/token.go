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

func (controller *tokenController) Create(w http.ResponseWriter, r *http.Request) {
	payload := readNewBody(r)
	generatedToken, err := utils.GenerateToken()
	if err != nil {
		fmt.Errorf("Cannot generate a new token : %s", err.Error())

		answerPayload := &utils.HTTPError{
			Msg:  "Cannot generate a new token",
			Code: "TOKEN01",
		}

		w.WriteHeader(http.StatusInternalServerError)
		answer, _ := json.Marshal(answerPayload)
		w.Write(answer)
		return
	}

	err = controller.DB.C("tokens").Insert(&models.Token{
		URL:   payload.URL,
		Token: generatedToken,
	})
	if err != nil {
		fmt.Errorf("Cannot save the new token in the database : %s", err.Error())
	}

	answerPayload := &createResponsePayload{
		Token: generatedToken,
		URL:   payload.URL,
	}
	answer, _ := json.Marshal(answerPayload)
	w.Write(answer)
}

// createRequestPayload should handle the request body
type createRequestPayload struct {
	URL string `json:"url"`
}

// createResponsePayload used to marshal the response
type createResponsePayload struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

func readNewBody(r *http.Request) *createRequestPayload {
	bodyBuf := make([]byte, r.ContentLength)
	r.Body.Read(bodyBuf)

	payload := new(createRequestPayload)
	err := json.Unmarshal(bodyBuf, payload)
	if err != nil {
		fmt.Errorf("Cannot unmarshal the request body : %s", err.Error())
	}

	return payload
}

// === Index Token ===

type indexResponsePayload struct {
	Data []models.Token
}

func (controller *tokenController) Index(w http.ResponseWriter, r *http.Request) {
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

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Redirection not found"))
		return
	}

	if len(tokenModel.URL) == 0 {
		fmt.Println("Token not found")

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Redirection not found"))
		return
	}

	http.Redirect(w, r, tokenModel.URL, http.StatusPermanentRedirect)
}
