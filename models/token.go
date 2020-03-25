package models

// Token database model
type Token struct {
	URL   string `json:"url" bson:"url"`
	Token string `json:"token" bson:"token"`
}
