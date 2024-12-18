
# Auth REST API Service

This project implements an authentication REST API service with the following features:

- Sign-Up: Create a new user account using email and password.
- Sign-In: Authenticate user credentials and return a JWT token.
- Revoke Token: Revoke an active JWT token.
- Refresh Token: Refresh an expired or about-to-expire JWT token.
- Token Validation: Middleware to protect routes using JWT.

## Tech Stack

Backend - Golang

Database - NoSql(MongoDB)

Caching - Redis

## Features Implemented

- Sign-Up: User registration with hashed password storage.
- Sign-In: JWT-based authentication.
- Revoke Token: Invalidate an existing token via Redis.
- Refresh Token: Generate a new token before expiry.
- Token Validation: Middleware to protect routes using JWT.

## Steps to Run

- git clone https://github.com/VARUNSHARMA55665/Auth-Rest-Api.git

- docker-compose up --build

## Curl
1. SignUp :
curl -X 'POST' \
  'http://localhost:8080/api/auth-rest-api/user/signUp' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "emailId": "testemail@gmail.com",
  "password": "pass@secret"
}'

2 SignIn :
curl -X 'POST' \
  'http://localhost:8080/api/auth-rest-api/user/signIn' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "emailId": "testemail@gmail.com",
  "password": "pass@secret"
}'

3 RefreshToken (provide valid authtoken) :
curl -X 'POST' \
  'http://localhost:8080/api/auth-rest-api/user/auth/refreshToken' \
  -H 'accept: application/json' \
  -H 'P-DeviceType: web' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzUxMTIzMzYsImlhdCI6MTczNDUwNzUzNiwic3ViamVjdCI6ImNoZWNrMzJAZ21haWwuY29tIn0.cEnSrV40XV55H45bhtkMT8NrpaK2OXkK4WWtgCtJyA9DP6ueFyYQYG48xYiWw6QRBxFwjGCxpw6ZkzITXm5wDA' \
  -d ''

4 RevokeToken (provide valid authtoken) :
curl -X 'POST' \
  'http://localhost:8080/api/auth-rest-api/user/auth/revokeToken' \
  -H 'accept: application/json' \
  -H 'P-DeviceType: web' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzUxMTIzNjIsImlhdCI6MTczNDUwNzU2Miwic3ViamVjdCI6ImNoZWNrMzJAZ21haWwuY29tIn0.05sb2CqIUPj1OLwf_aidvv_VquvWGgOHKM9f8Y-Vv5gpHL9P6iOS9SV7IJAMizdKRBqvWJjlQwdBoii_GK7oxw' \
  -d ''