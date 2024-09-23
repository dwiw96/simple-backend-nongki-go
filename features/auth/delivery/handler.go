package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	auth "simple-backend-nongki-go/features/auth"
	middleware "simple-backend-nongki-go/middleware"
	responses "simple-backend-nongki-go/utils/responses"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
)

type authDelivery struct {
	router   *httprouter.Router
	service  auth.ServiceInterface
	validate *validator.Validate
	trans    ut.Translator
}

func NewAuthDelivery(router *httprouter.Router, service auth.ServiceInterface, pool *pgxpool.Pool, client *redis.Client, ctx context.Context) {
	handler := &authDelivery{
		router:   router,
		service:  service,
		validate: validator.New(),
	}

	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(handler.validate, trans)
	handler.trans = trans

	router.POST("/api/signup", handler.SignUp)
	router.POST("/api/login", handler.LogIn)
	router.POST("/api/logout", middleware.AuthMiddleware(ctx, pool, client, handler.LogOut))
}

func translateError(trans ut.Translator, err error) (errTrans []string) {
	errs := err.(validator.ValidationErrors)
	a := (errs.Translate(trans))
	for _, val := range a {
		errTrans = append(errTrans, val)
	}

	return
}

func (d *authDelivery) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request signupRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		responses.ErrorJSON(w, http.StatusUnprocessableEntity, err.Error(), r.RemoteAddr)
		return
	}

	err = d.validate.Struct(request)
	if err != nil {
		errTranslated := translateError(d.trans, err)
		responses.ErrorJSON(w, 422, errTranslated, r.RemoteAddr)
		return
	}

	signupInput := toSignUpRequest(request)
	user, code, err := d.service.SignUp(signupInput)
	if err != nil {
		responses.ErrorJSON(w, code, err.Error(), r.RemoteAddr)
		return
	}

	response := toSignUpResponse(user)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(responses.SuccessWithDataResponse(response, 201, "SignUp success"))
}

func (d *authDelivery) LogIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request signinRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		responses.ErrorJSON(w, http.StatusUnprocessableEntity, err.Error(), r.RemoteAddr)
		return
	}

	err = d.validate.Struct(request)
	if err != nil {
		errTranslated := translateError(d.trans, err)
		responses.ErrorJSON(w, 422, errTranslated, r.RemoteAddr)
		return
	}

	user, token, code, err := d.service.LogIn(auth.LoginRequest(request))
	if err != nil {
		responses.ErrorJSON(w, code, err.Error(), r.RemoteAddr)
		return
	}

	response := toSignUpResponse(user)

	w.Header().Set("Authorization", token)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(responses.SuccessWithDataResponse(response, 200, "Login success"))
}

func (d *authDelivery) LogOut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authPayload := r.Context().Value(middleware.PayloadKey).(*auth.JwtPayload)
	fmt.Println("handler, authPayload:", authPayload)

	err := d.service.LogOut(*authPayload)
	if err != nil {
		responses.ErrorJSON(w, 401, err.Error(), r.RemoteAddr)
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode("logout success")
}
