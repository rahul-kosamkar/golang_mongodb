package main

// ErrorResponse is struct for sending error message with code.
type ErrorResponse struct {
	Code    int
	Message string
}

// SuccessResponse is struct for sending success message with code.
type SuccessResponse struct {
	Code     int
	Message  string
}

// RegistationParams is struct to read the request body
type RegistationParams struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserDetails is struct used for user details
type UserDetails struct {
	Name     string
	Email    string
	Password string
}

//ResetPasswordParams is struct used for reset password
type ResetPasswordParams struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

