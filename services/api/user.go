package services

import (
	apihelpers "auth-rest-api/apiHelpers"
	"auth-rest-api/models"
	"net/http"
)

type UserObj struct{}

func InitUser() UserObj {
	userObj := UserObj{}
	return userObj
}

func (obj UserObj) UserSignUp(req models.LogInReq) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	apiRes.Data = nil
	apiRes.Status = true
	apiRes.Message = "SUCCESS"
	return http.StatusOK, apiRes
}

func (obj UserObj) UserSignIn(req models.LogInReq) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	apiRes.Data = nil
	apiRes.Status = true
	apiRes.Message = "SUCCESS"
	return http.StatusOK, apiRes
}
