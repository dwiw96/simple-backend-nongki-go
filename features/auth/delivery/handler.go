package delivery

import (
	"encoding/json"
	"net/http"

	auth "simple-backend-nongki-go/features/auth"
	responses "simple-backend-nongki-go/utils/responses"

	"github.com/julienschmidt/httprouter"
)

type authDelivery struct {
	router  *httprouter.Router
	service auth.ServiceInterface
}

func NewAuthDelivery(router *httprouter.Router, serive auth.ServiceInterface) {
	handler := &authDelivery{
		router:  router,
		service: serive,
	}

	router.POST("/api/signup", handler.SignUp)
	router.POST("/api/signin", handler.SignIn)
}

func (d *authDelivery) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request signupRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		responses.ErrorJSON(w, http.StatusUnprocessableEntity, err.Error(), r.RemoteAddr)
	}

	signupInput := toSignUpRequest(request)
	user, code, err := d.service.SignUp(signupInput)
	if err != nil {
		responses.ErrorJSON(w, code, err.Error(), r.RemoteAddr)
	}

	response := toSignUpResponse(*user)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	responses.SuccessWithDataResponse(response, "signup success")
}

func (d *authDelivery) SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request signinRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		responses.ErrorJSON(w, http.StatusUnprocessableEntity, err.Error(), r.RemoteAddr)
	}

	user, token, code, err := d.service.SignIn(auth.SigninRequest(request))
	if err != nil {
		responses.ErrorJSON(w, code, err.Error(), r.RemoteAddr)
	}

	response := toSignUpResponse(*user)

	w.Header().Set("Authorization", token)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	responses.SuccessWithDataResponse(response, "signup success")
}
