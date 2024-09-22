package middleware

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	// "net/http"

	auth "simple-backend-nongki-go/features/auth"
	// "github.com/julienschmidt/httprouter"
)

type KeyCache struct {
	key        *rsa.PrivateKey
	expiration time.Time
	mutex      sync.Mutex
}

var keyCache = &KeyCache{
	key:        nil,
	expiration: time.Now(),
}

// func TokenMiddleware(db *sql.DB, next httprouter.Handle) httprouter.Handle {
// 	return (func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 		if r.RequestURI == "/api/signin" {
// 			next(w, r, ps)
// 			return
// 		}

// 		key, err := LoadKey(db)
// 		if err != nil {
// 			log.Println(err)
// 			return err
// 		}

// 		tokenString, err := GetTokenHeader(c)
// 		if err != nil {
// 			log.Println(err)
// 			return err
// 		}

// 		legitimate := Verify(tokenString, key)
// 		if !legitimate {
// 			log.Println("token illegitimate")
// 			return errors.New("wrong token")
// 		}

// 		tokenData := ReadToken(tokenString)

// 		newToken := GetUpdatedToken(tokenData, key)

// 		w.Header().Set("Access-Control-Expose-Headers", "Authorization,Access-Control-Allow-Origin,Access-Control-Allow-Credentials,Access-Control-Allow-Methods,Access-Control-Allow-Headers")
// 		w.Header().Set("Authorization", newToken)
// 		next(w, r, ps)
// 	})
// }

func LoadKey(db *sql.DB) (key *rsa.PrivateKey, err error) {
	if keyCache.key != nil && time.Now().Before(keyCache.expiration) {
		return keyCache.key, nil
	}

	query := "select private_key from sec_m"
	var keyBytes []byte
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&keyBytes)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		log.Println("Load success")

		privateKey, err := x509.ParsePKCS1PrivateKey(keyBytes)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		keyCache.mutex.Lock()
		defer keyCache.mutex.Unlock()
		keyCache.key = key
		keyCache.expiration = time.Now().Add(time.Hour)

		return privateKey, nil
	}

	return nil, errors.New("no private key found in database")
}

// func GetTokenHeader(c echo.Context) (token string, err error) {
// 	authHeader := c.Request().Header.Get("Authorization")
// 	if authHeader == "" {
// 		return "", errors.New("no authorization header found")
// 	}

// 	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
// 	return tokenString, nil
// }

// func CheckAuthorization(c echo.Context, key *rsa.PrivateKey) (tokenData auth.JwtPayload, newTokenString string, err error) {
// 	tokenString, err := GetTokenHeader(c)
// 	if err != nil {
// 		return auth.JwtPayload{}, "", err
// 	}

// 	legitimate := Verify(tokenString, key)
// 	if !legitimate {
// 		return auth.JwtPayload{}, "", errors.New("wrong token")
// 	}

// 	tokenData = ReadToken(tokenString)

// 	if time.Now().UTC().Add(time.Hour * 9).After(tokenData.Exp.Add(time.Minute * 5)) {
// 		return auth.JwtPayload{}, "", errors.New("expired token")
// 	}

// 	if tokenData.Exp.Before(time.Now().UTC().Add(time.Hour * 9)) {
// 		newToken := GetUpdatedToken(tokenData, key)

// 		return tokenData, newToken, nil
// 	}

// 	return tokenData, "", nil
// }

func GetSignedToken(reqData auth.User, key *rsa.PrivateKey) string {
	privateKey := key
	header := auth.JwtHeader{"RS256", "JWT"}
	nowTime := time.Now().UTC().Add(time.Hour * 9).Add(time.Minute * 60)
	payload := auth.JwtPayload{
		UserID:  reqData.ID,
		Name:    reqData.FirstName + reqData.MiddleName + reqData.LastName,
		Email:   reqData.Email,
		Address: reqData.Address,
		Iat:     nowTime,
		Exp:     nowTime,
	}
	headerBytes, err := json.Marshal(header)
	payloadBytes, err := json.Marshal(payload)
	encodedHeader := base64.URLEncoding.EncodeToString(headerBytes)
	encodedPayload := base64.URLEncoding.EncodeToString(payloadBytes)
	headPay := encodedHeader + "." + encodedPayload
	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(headPay))
	if err != nil {
		log.Println(err)
		//panic(err)
	}
	msgHashSum := msgHash.Sum(nil)
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		log.Println(err)
		//panic(err)
	}
	encodedSignature := base64.URLEncoding.EncodeToString(signature)
	token := "Bearer " + headPay + "." + encodedSignature
	return token
}

func GetUpdatedToken(payload auth.JwtPayload, key *rsa.PrivateKey) string {
	privateKey := key
	header := auth.JwtHeader{"RS256", "JWT"}
	nowTime := time.Now().UTC().Add(time.Hour * 9)
	payload.Exp = nowTime.Add(time.Minute * 15)
	headerBytes, err := json.Marshal(header)
	payloadBytes, err := json.Marshal(payload)
	encodedHeader := base64.URLEncoding.EncodeToString(headerBytes)
	encodedPayload := base64.URLEncoding.EncodeToString(payloadBytes)
	headPay := encodedHeader + "." + encodedPayload
	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(headPay))
	if err != nil {
		panic(err)
	}
	msgHashSum := msgHash.Sum(nil)
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		panic(err)
	}
	encodedSignature := base64.URLEncoding.EncodeToString(signature)
	token := "Bearer " + headPay + "." + encodedSignature
	return token
}

func Verify(token string, key *rsa.PrivateKey) bool {
	if len(token) < 20 {
		return false
	}
	spaceLastIndex := strings.LastIndex(token, " ")
	if len(token) < spaceLastIndex+2 {

	} else {
		token = token[spaceLastIndex+1:]
	}

	lastCommaIndex := strings.LastIndex(token, ".")
	if len(token) < lastCommaIndex+2 {
		return false
	}

	signature := token[lastCommaIndex+1:]
	verifyPart := token[:lastCommaIndex]

	decodedSignature, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		log.Println("could not decode signature: ", err)
		return false
	}

	// Verify the signature using the same algorithm and method as used in the GetUpdatedToken function
	hash := sha256.New()
	hash.Write([]byte(verifyPart))
	hashed := hash.Sum(nil)
	opts := rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA256}
	err = rsa.VerifyPSS(&key.PublicKey, crypto.SHA256, hashed, decodedSignature, &opts)
	if err != nil {
		log.Println("could not verify signature: ", err)
		return false
	}

	log.Println("signature verified")
	return true
}

func VerifyExpired() bool {

	return true
}

func ReadToken(token string) auth.JwtPayload {
	payloadModel := &auth.JwtPayload{}
	tokenArr := strings.Split(token, ".")
	if len(tokenArr) != 3 {
		log.Println("token token!")
		return *payloadModel
	}
	payload := tokenArr[1]

	payloadBytes, err := base64.URLEncoding.DecodeString(payload)

	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(payloadBytes, payloadModel)

	return *payloadModel
}
