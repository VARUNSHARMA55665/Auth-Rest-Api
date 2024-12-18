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
	"golang.org/x/crypto/bcrypt"
)

type UserObj struct{}

func InitUser() UserObj {
	userObj := UserObj{}
	return userObj
}

func (obj UserObj) UserSignUp(req models.LogInReq) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	// Initialize the API response
	apiRes.Data = nil
	apiRes.Status = false

	// Check if a user already exists with the given email
	var existingUser models.MongoSignup
	err := db.FindOneMongo(constants.CLIENTCOLLECTION, bson.M{"emailid": req.EmailId}, &existingUser)
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

	// Prepare data for MongoDB
	var newUser models.MongoSignup
	newUser.EmailId = req.EmailId
	hashPassword, err := hashPassword(req.Passwrod)
	if err != nil {
		log.Println("UserSignUp error in hashing the password", err)
		return apihelpers.SendInternalServerError()
	}

	newUser.Passwrod = hashPassword
	newUser.CreatedAt = time.Now().Unix()
	newUser.UpdatedAt = time.Now().Unix()

	// Upsert the user data into the database
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
	err = helpers.SetRedis(redisAuthKey, jwtToken, constants.TokenTTL) // set token for one day
	if err != nil {
		log.Println(" VerifyOtp SetRedis failed, err = ", err, " emailId:", req.EmailId)
		return apihelpers.SendInternalServerError()
	}

	var signupRes models.LogInRes
	signupRes.Authorization = "Bearer " + jwtToken

	// Success response
	log.Println("UserSignUp: Successfully signed up user with emailId:", req.EmailId)

	apiRes.Data = signupRes
	apiRes.Status = true
	apiRes.Message = "User signed up successfully"

	return http.StatusOK, apiRes
}

func (obj UserObj) UserSignIn(req models.LogInReq) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	// Initialize the API response
	apiRes.Data = nil
	apiRes.Status = false

	// Fetch user from the database
	var existingUser models.MongoSignup
	err := db.FindOneMongo(constants.CLIENTCOLLECTION, bson.M{"emailid": req.EmailId}, &existingUser)
	if err != nil {
		if err.Error() == constants.MongoNoDocError {
			// User not found
			log.Println("UserSignIn: User not found with emailId:", req.EmailId)
			apiRes.Message = "Invalid email or password"
			return http.StatusUnauthorized, apiRes
		}
		// Database error
		log.Println("UserSignIn error in FindOneMongo err:", err, " emailId:", req.EmailId)
		return apihelpers.SendInternalServerError()
	}

	err = checkPasswordHash(req.Passwrod, existingUser.Passwrod)
	// Validate password
	if err != nil { // existingUser.Passwrod != req.Passwrod { // In production, compare hashed passwords
		log.Println("UserSignIn: Invalid password for emailId:", req.EmailId, " err: ", err)
		apiRes.Message = "Invalid email or password"
		return http.StatusUnauthorized, apiRes
	}

	// Generate JWT token
	jwtToken, err := helpers.GenerateJWT(req.EmailId)
	if err != nil {
		log.Println("UserSignIn GenerateJWT failed, err =", err, " emailId:", req.EmailId)
		return apihelpers.SendInternalServerError()
	}

	// Store token in Redis
	redisAuthKey := "auth|" + req.EmailId
	err = helpers.SetRedis(redisAuthKey, jwtToken, constants.TokenTTL) // Set token for one day
	if err != nil {
		log.Println("UserSignIn SetRedis failed, err =", err, " emailId:", req.EmailId)
		return apihelpers.SendInternalServerError()
	}

	// Prepare Response
	var signinRes models.LogInRes
	signinRes.Authorization = "Bearer " + jwtToken

	apiRes.Data = signinRes
	apiRes.Status = true
	apiRes.Message = "User signed in successfully"

	// Success Log
	log.Println("UserSignIn: Successfully signed in user with emailId:", req.EmailId)

	return http.StatusOK, apiRes
}

// hashPassword hashes a plain text password using bcrypt
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func checkPasswordHash(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

func (obj UserObj) RevokeToken(authHeader string) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	claims, err := helpers.ValidateToken(authHeader[7:]) // size check on authHeader is alreay added in middleware
	if err != nil {
		apiRes.Data = nil
		apiRes.Status = false
		apiRes.Message = "Invalid or expired token"
		return http.StatusUnauthorized, apiRes
	}

	// Get user email from the claims
	email := string(claims)
	// Revoke Token by deleting from Redis
	redisAuthKey := "auth|" + email
	err = helpers.DelRedis(redisAuthKey).Err()
	if err != nil {
		log.Println("RevokeToken Failed to delete token from Redis, err =", err)
		return apihelpers.SendInternalServerError()
	}

	log.Println("RevokeToken Token revoked for email:", email)

	// Success Response
	apiRes.Status = true
	apiRes.Message = "Token revoked successfully"
	return http.StatusOK, apiRes
}

func (obj UserObj) RefreshToken(authHeader string) (int, apihelpers.APIRes) {
	var apiRes apihelpers.APIRes

	// Initialize response
	apiRes.Data = nil
	apiRes.Status = false
	apiRes.Message = "Token refresh failed"

	// Verify the current JWT Token
	claims, err := helpers.ValidateToken(authHeader[7:]) // size check on authHeader is alreay added in middleware
	if err != nil {
		apiRes.Data = nil
		apiRes.Status = false
		apiRes.Message = "Invalid or expired token"
		return http.StatusUnauthorized, apiRes
	}

	// Get user email from the claims
	email := string(claims)
	// Generate a new JWT token
	newToken, err := helpers.GenerateJWT(email)
	if err != nil {
		log.Println("RefreshToken: Failed to generate new token, err =", err)
		return apihelpers.SendInternalServerError()
	}

	redisAuthKey := "auth|" + email

	err = helpers.SetRedis(redisAuthKey, newToken, constants.TokenTTL) // Set token TTL to 1 day
	if err != nil {
		log.Println("RefreshToken: Failed to update Redis with new token, err =", err)
		return apihelpers.SendInternalServerError()
	}

	var refreshRes models.LogInRes
	refreshRes.Authorization = "Bearer " + newToken

	log.Println("RefreshToken: New token issued for email:", email)

	// Prepare the Response
	apiRes.Data = refreshRes
	apiRes.Status = true
	apiRes.Message = "Token refreshed successfully"

	return http.StatusOK, apiRes
}
