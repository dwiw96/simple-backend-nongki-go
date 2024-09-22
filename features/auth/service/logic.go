package service

import (
	"fmt"
	auth "simple-backend-nongki-go/features/auth"
	middleware "simple-backend-nongki-go/middleware"
	password "simple-backend-nongki-go/utils/password"
)

type authService struct {
	repo auth.RepositoryInterface
}

func NewAuthService(repo auth.RepositoryInterface) auth.ServiceInterface {
	return &authService{
		repo: repo,
	}
}

func (s *authService) SignUp(input auth.SignupRequest) (user *auth.User, code int, err error) {
	// check if email is registered
	// --> true, send error 409 conflict
	resCheckEmail, err := s.repo.CheckEmail(input.Email)
	if err != nil {
		return nil, 500, err
	}

	if resCheckEmail != 0 {
		return nil, 409, fmt.Errorf("email is registered")
	}

	// hashing password
	input.HashedPassword, err = password.HashingPassword(input.HashedPassword)
	if err != nil {
		return nil, 500, err
	}

	// insert new user
	user, err = s.repo.InsertUser(input)
	if err != nil {
		return nil, 400, err
	}

	return user, 0, nil
}

func (s *authService) SignIn(input auth.SigninRequest) (user *auth.User, token string, code int, err error) {
	user, err = s.repo.ReadUser(input.Email)
	if err != nil {
		return nil, "", 400, err
	}

	err = password.VerifyHashPassword(input.Password, user.HashedPassword)
	if err != nil {
		return nil, "", 400, err
	}

	key, err := s.repo.LoadKey()
	if err != nil {
		return nil, "", 500, fmt.Errorf("load key error: %w", err)
	}

	token = middleware.GetSignedToken(*user, key)

	return user, token, 200, nil
}
