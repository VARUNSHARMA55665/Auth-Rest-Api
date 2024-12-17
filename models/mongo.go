package models

type MongoSignup struct {
	EmailId   string `json:"emailId"`
	Passwrod  string `json:"password"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}
