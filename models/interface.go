package models

import apihelpers "auth-rest-api/apiHelpers"

type UserProvider interface {
	UserSignUp(req LogInReq) (int, apihelpers.APIRes)
	UserSignIn(req LogInReq) (int, apihelpers.APIRes)
	RevokeToken(token string) (int, apihelpers.APIRes)
	RefreshToken(token string) (int, apihelpers.APIRes)
}
