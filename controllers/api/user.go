package controllers

import (
	apihelpers "auth-rest-api/apiHelpers"
	"auth-rest-api/constants"
	"auth-rest-api/models"
	"encoding/json"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

var theUserProvider models.UserProvider

func InitUserProvider(provider models.UserProvider) {
	theUserProvider = provider
}

// Signup
// @Description enter email and password for signup
// @Tags users V1
// @Accept json
// @Produce json
// @Param request body models.LogInReq true "Signup"
// @Success 200 {object} apihelpers.APIRes{data=models.LogInRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 500 {object} apihelpers.APIRes
// @Router /api/auth-rest-api/user/signUp [post]
func SignUp(c *gin.Context) {
	var userDetails models.LogInReq

	err := json.NewDecoder(c.Request.Body).Decode(&userDetails)
	if err != nil {
		log.Println("Signup (Controllers) error:", err)
		apihelpers.ErrorMessageController(c, constants.InvalidRequest)
		return
	}

	// Normalize the email (to lower case)
	userDetails.EmailId = strings.ToLower(userDetails.EmailId)

	log.Println("Signup (controller), reqParams email:", userDetails.EmailId)

	code, resp := theUserProvider.UserSignUp(userDetails)
	apihelpers.CustomResponse(c, code, resp)

}

// SignIn
// @Description enter email and password for signin
// @Tags users V1
// @Accept json
// @Produce json
// @Param request body models.LogInReq true "signin"
// @Success 200 {object} apihelpers.APIRes{data=models.LogInRes}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 500 {object} apihelpers.APIRes
// @Router /api/auth-rest-api/user/signIn [post]
func SignIn(c *gin.Context) {
	var userDetails models.LogInReq

	err := json.NewDecoder(c.Request.Body).Decode(&userDetails)
	if err != nil {
		log.Println("SignIn (Controllers) error:", err)
		apihelpers.ErrorMessageController(c, constants.InvalidRequest)
		return
	}

	userDetails.EmailId = strings.ToLower(userDetails.EmailId)

	log.Println("SignIn (controller), reqParams email:", userDetails.EmailId)

	code, resp := theUserProvider.UserSignIn(userDetails)
	apihelpers.CustomResponse(c, code, resp)

}

// RevokeToken
// @Description Revoke Token
// @Tags users V1
// @Accept json
// @Produce json
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param Authorization header string true "Authorization Header"
// @Success 200 {object} apihelpers.APIRes
// @Failure 400 {object} apihelpers.APIRes
// @Failure 500 {object} apihelpers.APIRes
// @Router /api/auth-rest-api/user/auth/revokeToken [post]
func RevokeToken(c *gin.Context) {
	var requestH models.ReqHeader
	if err := c.ShouldBindHeader(&requestH); err != nil {
		log.Println("RevokeToken error in ShouldBindHeader err: ", err)
		apihelpers.SendInternalServerError()
		return
	}

	log.Println("In RevokeToken (controller)")

	code, resp := theUserProvider.RevokeToken(requestH.Authorization)
	apihelpers.CustomResponse(c, code, resp)

}

// RefreshToken
// @Description Refresh Token
// @Tags users V1
// @Accept json
// @Produce json
// @Param P-DeviceType header string true "P-DeviceType Header"
// @Param Authorization header string true "Authorization Header"
// @Success 200 {object} apihelpers.APIRes{data=models.AuthToken}
// @Failure 400 {object} apihelpers.APIRes
// @Failure 500 {object} apihelpers.APIRes
// @Router /api/auth-rest-api/user/auth/refreshToken [post]
func RefreshToken(c *gin.Context) {
	var requestH models.ReqHeader
	if err := c.ShouldBindHeader(&requestH); err != nil {
		log.Println("RevokeToken error in ShouldBindHeader err: ", err)
		apihelpers.SendInternalServerError()
		return
	}

	log.Println("In RefreshToken (controller)")

	code, resp := theUserProvider.RefreshToken(requestH.Authorization)
	apihelpers.CustomResponse(c, code, resp)

}
