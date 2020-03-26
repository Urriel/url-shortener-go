package models

// Token database model
type Token struct {
	URL    string `json:"url" bson:"url"`
	Token  string `json:"token" bson:"token"`
	Visits uint32 `json:"visits" bson:"visits"`
}
