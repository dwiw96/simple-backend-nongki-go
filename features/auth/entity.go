package auth

import (
	"crypto/rsa"
	"time"
)

type SignupRequest struct {
	FirstName      string `json:"first_name"`
	MiddleName     string `json:"middle_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Address        string `json:"address"`
	Gender         string `json:"gender"`
	MaritalStatus  string `json:"marital_status"`
	HashedPassword string `json:"password"`
}

type User struct {
	ID             int64
	Fullname       string
	FirstName      string `json:"first_name"`
	MiddleName     string `json:"middle_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Address        string `json:"address"`
	Gender         string `json:"gender"`
	MaritalStatus  string `json:"marital_status"`
	HashedPassword string `json:"hashed_password"`
	CreatedAt      time.Time
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JwtPayload struct {
	UserID  int64     `json:"usi"`
	Iss     string    `json:"iss"`
	Name    string    `json:"nam"`
	Email   string    `json:"eml"`
	Address string    `json:"adr,omitempty"`
	Iat     time.Time `json:"iat"`
	Exp     time.Time `json:"exp"`
}

type RepositoryInterface interface {
	CheckEmail(email string) (result int, err error)
	ReadUser(email string) (result *User, err error)
	InsertUser(input SignupRequest) (result *User, err error)
	LoadKey() (key *rsa.PrivateKey, err error)
}

type ServiceInterface interface {
	SignUp(input SignupRequest) (user *User, code int, err error)
	SignIn(input SigninRequest) (user *User, token string, code int, err error)
}
