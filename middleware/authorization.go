package middleware

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	// "net/http"

	auth "simple-backend-nongki-go/features/auth"
	// "github.com/julienschmidt/httprouter"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetToken(reqData auth.User, key *rsa.PrivateKey) (token string, err error) {
	nowTime := time.Now().UTC()
	expTime := nowTime.Add(time.Hour * 1)

	t := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"iss":     "nongki",
			"userID":  reqData.ID,
			"name":    reqData.Fullname,
			"email":   reqData.Email,
			"address": reqData.Address,
			"iat":     nowTime.Unix(),
			"exp":     expTime.Unix(),
		})

	token, err = t.SignedString(key)

	return
}

func VerifyToken(userToken string, key *rsa.PrivateKey) (bool, error) {
	jwtToken, err := jwt.Parse(userToken, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"])
		}

		return &key.PublicKey, nil
	})

	if err != nil {
		log.Println("err:", err)
		return false, err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		if claims["iss"] != "nongki" {
			errMsg := fmt.Errorf("claim iss is wrong")
			log.Println(errMsg)
			return false, errMsg
		}
	} else {
		errMsg := fmt.Errorf("iss payload not found")
		log.Println(errMsg)
		return false, errMsg
	}

	return true, nil
}

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

func LoadKey(ctx context.Context, conn *pgxpool.Pool) (key *rsa.PrivateKey, err error) {
	if keyCache.key != nil && time.Now().Before(keyCache.expiration) {
		return keyCache.key, nil
	}

	query := "select private_key from sec_m"
	var keyBytes []byte
	rows, err := conn.Query(ctx, query)
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

func GetUpdatedToken(payload auth.JwtPayload, key *rsa.PrivateKey) string {
	privateKey := key
	header := auth.JwtHeader{Alg: "RS256", Typ: "JWT"}
	nowTime := time.Now().UTC()
	payload.Exp = nowTime.Add(time.Minute * 60)
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
