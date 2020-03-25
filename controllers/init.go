package controllers

import (
	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2"
)

// InitControllers initialize all controllers
func InitControllers(mux *web.Mux, database *mgo.Database) {
	NewTokenController(mux, database)
}
