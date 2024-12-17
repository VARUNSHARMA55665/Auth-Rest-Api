package models

type LogInReq struct {
	EmailId  string `json:"emailId" example:"abc@gmail.com"`
	Passwrod string `json:"password"`
}

type LogInRes struct {
	Authorization string `header:"Authorization"`
}

type AuthToken struct {
	Token string `header:"token"`
}
