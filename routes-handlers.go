package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"unicode"

	"go.mongodb.org/mongo-driver/bson"
)

// SignUpUser Used for Registering the Users
func SignUpUser(response http.ResponseWriter, request *http.Request) {
	var registationRequest RegistationParams
	var result UserDetails
	var errorResponse = ErrorResponse {
		Code: http.StatusInternalServerError, Message: "It's not you it's me.",
	}

	decoder := json.NewDecoder(request.Body)
	decoderErr := decoder.Decode(&registationRequest)
	defer request.Body.Close()

	if decoderErr != nil {
		returnErrorResponse(response, request, errorResponse)
	} else {
		errorResponse.Code = http.StatusBadRequest
		if registationRequest.UserName == "" {
			errorResponse.Message = "username can't be empty"
			returnErrorResponse(response, request, errorResponse)
		} else if registationRequest.Email == "" {
			errorResponse.Message = "email can't be empty"
			returnErrorResponse(response, request, errorResponse)
		} else {
			isValid, errs := isValidPassword(registationRequest.Password, registationRequest.UserName, 8, 16)

			if (!isValid) {
				returnErrorResponse(response, request, errs)
				return
			}
			collection := Client.Database("test").Collection("users")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			var err = collection.FindOne(ctx, bson.M{
				"username": registationRequest.UserName,
			}).Decode(&result)
			defer cancel()

			if err == nil {
				errorResponse.Message = "Username already exists in the system, Please try with different username."
				returnErrorResponse(response, request, errorResponse)
				return
			}

			ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
			_, databaseErr := collection.InsertOne(ctx1, bson.M{
				"email": registationRequest.Email,
				"password": registationRequest.Password,
				"username": registationRequest.UserName,
			})
			defer cancel1()

			if databaseErr != nil {
				returnErrorResponse(response, request, errorResponse)
			}

			var successResponse = SuccessResponse {
				Code:     http.StatusOK,
				Message:  "Successfully registered.",
			}

			successJSONResponse, jsonError := json.Marshal(successResponse)

			if jsonError != nil {
				returnErrorResponse(response, request, errorResponse)
			}
			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(successResponse.Code)
			response.Write(successJSONResponse)
		}
	}
}

// ResetPassword Used for reseting user's password
func ResetPassword(response http.ResponseWriter, request *http.Request) {
	var resetPasswordRequest ResetPasswordParams
	var result UserDetails

	var errorResponse = ErrorResponse {
		Code: http.StatusInternalServerError, Message: "It's not you it's me.",
	}

	decoder := json.NewDecoder(request.Body)
	decoderErr := decoder.Decode(&resetPasswordRequest)
	defer request.Body.Close()

	if decoderErr != nil {
		returnErrorResponse(response, request, errorResponse)
	} else {
		errorResponse.Code = http.StatusBadRequest
		if resetPasswordRequest.UserName == "" {
			errorResponse.Message = "username can't be empty"
			returnErrorResponse(response, request, errorResponse)
		}
		isValid, errs := isValidPassword(resetPasswordRequest.Password, resetPasswordRequest.UserName, 8, 16)

		if (!isValid) {
			returnErrorResponse(response, request, errs)
			return
		}

		collection := Client.Database("test").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var err = collection.FindOne(ctx, bson.M{
			"username": resetPasswordRequest.UserName,
		}).Decode(&result)
		defer cancel()

		if err != nil {
			returnErrorResponse(response, request, errorResponse)
		}
		ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)

		filter := bson.M{
			"username": resetPasswordRequest.UserName,
		}
		update := bson.M{
			"$set": bson.M{
				"password": resetPasswordRequest.Password,
			},
		}
		_, databaseErr := collection.UpdateOne(ctx1, filter, update)
		defer cancel1()

		if databaseErr != nil {
			returnErrorResponse(response, request, errorResponse)
		}

		var successResponse = SuccessResponse {
			Code:     http.StatusOK,
			Message:  "Password updated successfully.",
		}

		successJSONResponse, jsonError := json.Marshal(successResponse)

		if jsonError != nil {
			returnErrorResponse(response, request, errorResponse)
		}
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(successResponse.Code)
		response.Write(successJSONResponse)

	}
}

//Function to validate password
func isValidPassword(password string, username string, min, max int) (isValid bool, err ErrorResponse) {
    var (
        isMin   bool
        special bool
        number  bool
        upper   bool
        lower   bool
    )
	var errorResponse = ErrorResponse{
		Code: http.StatusInternalServerError, Message: "It's not you it's me.",
	}

    //test for the muximum and minimum characters required for the password string
    if len(password) < min || len(password) > max {
        isMin = false
		errorResponse.Message = "length should be " + strconv.Itoa(min) + " to " + strconv.Itoa(max)
		return false, errorResponse;
    }

    for _, c := range password {
        // Optimize perf if all become true before reaching the end
        if special && number && upper && lower && isMin {
            break
        }

        // else go on switching
        switch {
        case unicode.IsUpper(c):
            upper = true
        case unicode.IsLower(c):
            lower = true
        case unicode.IsNumber(c):
            number = true
        case unicode.IsPunct(c) || unicode.IsSymbol(c):
            special = true
        }
    }

    // Add custom error messages
    if !special {
		errorResponse.Message = "password should contain at least a single special character"
		return false, errorResponse
    }
    if !number {
        errorResponse.Message = "password should contain at least a single digit"
		return false, errorResponse
    }
    if !lower {
        errorResponse.Message = "password should contain at least a single lowercase letter"
		return false, errorResponse
    }
    if !upper {
        errorResponse.Message = "password should contain at least single uppercase letter"
		return false, errorResponse
    }
	if (username == password) {
		errorResponse.Message = "password can't be same as username"
		return false, errorResponse
	}

    // everyting is right
    return true, errorResponse
}

func returnErrorResponse(response http.ResponseWriter, request *http.Request, errorMesage ErrorResponse) {
	httpResponse := &ErrorResponse{Code: errorMesage.Code, Message: errorMesage.Message}
	jsonResponse, err := json.Marshal(httpResponse)
	if err != nil {
		panic(err)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(errorMesage.Code)
	response.Write(jsonResponse)
}
