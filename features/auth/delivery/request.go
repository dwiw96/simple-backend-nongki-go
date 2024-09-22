package delivery

import (
	auth "simple-backend-nongki-go/features/auth"
)

type signupRequest struct {
	FirstName     string `json:"first_name"`
	MiddleName    string `json:"middle_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Address       string `json:"address"`
	Gender        string `json:"gender"`
	MaritalStatus string `json:"marital_status"`
	Password      string `json:"password"`
}

func toSignUpRequest(input signupRequest) auth.SignupRequest {
	return auth.SignupRequest{
		FirstName:      input.FirstName,
		MiddleName:     input.MiddleName,
		LastName:       input.LastName,
		Email:          input.Email,
		Address:        input.Address,
		Gender:         input.Gender,
		MaritalStatus:  input.MaritalStatus,
		HashedPassword: input.Password,
	}
}

type signinRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
