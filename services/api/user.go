package services

import (
	apihelpers "auth-rest-api/apiHelpers"
	"auth-rest-api/constants"
	"auth-rest-api/db"
	"auth-rest-api/helpers"
	"auth-rest-api/models"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserObj struct{}

func InitUser() UserObj {
	userObj := UserObj{}
	return userObj
}

// func (obj UserObj) UserSignUp(req models.LogInReq) (int, apihelpers.APIRes) {
// 	var apiRes apihelpers.APIRes

// 	apiRes.Data = nil
// 	apiRes.Status = true
// 	apiRes.Message = "SUCCESS"
// 	return http.StatusOK, apiRes
// }

func (obj UserObj) UserSignUp(req models.LogInReq) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	// Initialize the API response
	apiRes.Data = nil
	apiRes.Status = false
	apiRes.Message = "Signup Failed"

	// Step 2: Check if a user already exists with the given email
	var existingUser models.LogInReq
	err := db.FindOneMongo(constants.CLIENTCOLLECTION, bson.M{"emailId": req.EmailId}, &existingUser)
	if err != nil && err.Error() != constants.MongoNoDocError {
		log.Println("UserSignUp error in findone mongo err: ", err, " emailId: ", req.EmailId)
		return apihelpers.SendInternalServerError()
	}
	if err == nil {
		// If no error, user already exists
		log.Println("UserSignUp: User already exists with emailId:", req.EmailId)
		apiRes.Message = "User already exists"
		return http.StatusBadRequest, apiRes
	}

	// Step 3: Prepare data for MongoDB
	var newUser models.MongoSignup
	newUser.EmailId = req.EmailId
	newUser.Passwrod = req.Passwrod
	newUser.CreatedAt = time.Now().Unix()
	newUser.UpdatedAt = time.Now().Unix()

	// Step 4: Upsert the user data into the database
	filter := bson.M{"emailid": req.EmailId}
	update := bson.M{"$set": newUser}
	opts := options.Update().SetUpsert(true)

	err = db.UpdateOneMongo(constants.CLIENTCOLLECTION, filter, update, opts)
	if err != nil {
		// Handle database error
		log.Println("UserSignUp: MongoDB Upsert failed, err =", err, " emailId:", req.EmailId)
		return apihelpers.SendInternalServerError()
	}

	// generate jwt token
	jwtToken, err := helpers.GenerateJWT(req.EmailId)
	if err != nil {
		log.Println("UserSignUp GenerateJWT failed, err = ", err, " emailId:", req.EmailId)
		return apihelpers.SendInternalServerError()
	}

	redisAuthKey := "auth|" + req.EmailId
	err = helpers.SetRedis(redisAuthKey, jwtToken, 24*60) // set token for one day
	if err != nil {
		log.Println(" VerifyOtp SetRedis failed, err = ", err, " emailId:", req.EmailId)
		return apihelpers.SendInternalServerError()
	}

	var signupRes models.LogInRes
	signupRes.Authorization = "Bearer " + redisAuthKey

	// Step 5: Success response
	log.Println("UserSignUp: Successfully signed up user with emailId:", req.EmailId)

	apiRes.Data = signupRes
	apiRes.Status = true
	apiRes.Message = "User signed up successfully"

	return http.StatusOK, apiRes
}

func (obj UserObj) UserSignIn(req models.LogInReq) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	apiRes.Data = nil
	apiRes.Status = true
	apiRes.Message = "SUCCESS"
	return http.StatusOK, apiRes
}
