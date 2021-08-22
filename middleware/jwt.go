package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwtMid "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	datalog "github.com/mmuflih/go-text-log"
	"github.com/mmuflih/golib/response"
)

type httpFunc func(http.ResponseWriter, *http.Request)

var jwtMiddleware *jwtMid.JWTMiddleware
var signingKey []byte
var myRole map[string][]string
var writeLog datalog.DataLog

func InitJWTMiddleware(secret []byte) {
	InitJWTMiddlewareCustomSigningKey(secret, jwt.SigningMethodES512)
}

func InitJWTMiddlewareCustomSigningKey(secret []byte, signingMethod jwt.SigningMethod) {
	writeLog = datalog.New("jwt-mid.log", true)
	signingKey = secret
	jwtMiddleware = jwtMid.New(jwtMid.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		},
		SigningMethod: signingMethod,
	})
}

func InitJWTMiddlewareWithRole(secret []byte, signingMethod jwt.SigningMethod, role map[string][]string) {
	signingKey = secret
	myRole = role
	InitJWTMiddlewareCustomSigningKey(secret, signingMethod)
}

func ExtractClaim(r *http.Request, key string) (interface{}, error) {
	tokenStr, err := jwtMiddleware.Options.Extractor(r)
	if err != nil {
		return "", err
	}

	if tokenStr == "" {
		return nil, errors.New("Request Unauthorized")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return signingKey, nil
	})
	if token == nil {
		return nil, errors.New("Request Unauthorized #2")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if claims[key] != nil {
			return claims[key], nil
		}
		return nil, errors.New("Claim data not found")
	}

	if !ok {
		return nil, errors.New("Token broken")
	}
	if !token.Valid {
		return nil, errors.New("Token invalid")
	}
	return nil, nil
}

func JWTMid(h httpFunc) httpFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := checkJWT(w, r, "")
		if err != nil {
			return
		}
		h(w, r)
	}
}

func JWTMidWithRole(h httpFunc, role string) httpFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := checkJWT(w, r, role)
		if err != nil {
			return
		}
		h(w, r)
	}
}

func checkJWT(w http.ResponseWriter, r *http.Request, role string) error {

	if !jwtMiddleware.Options.EnableAuthOnOptions {
		if r.Method == "OPTIONS" {
			return nil
		}
	}

	token, err := jwtMiddleware.Options.Extractor(r)
	if err != nil {
		eExtractor := errors.New("400")
		response.Exception(w, eExtractor, 400)
		return eExtractor
	}

	if token == "" {

		if jwtMiddleware.Options.CredentialsOptional {
			return nil
		}

		eReqiredToken := errors.New("Required authorization token not found")
		response.Exception(w, eReqiredToken, 401)
		return eReqiredToken
	}

	parsedToken, err := jwt.Parse(token, jwtMiddleware.Options.ValidationKeyGetter)
	if err != nil {
		ePassingToken := errors.New("Error parsing token: " + err.Error())
		response.Exception(w, ePassingToken, 401)
		return ePassingToken
	}

	if jwtMiddleware.Options.SigningMethod != nil && jwtMiddleware.Options.SigningMethod.Alg() != parsedToken.Header["alg"] {
		errorMsg := fmt.Sprintf("Expected %s signing method but token specified %s",
			jwtMiddleware.Options.SigningMethod.Alg(),
			parsedToken.Header["alg"])
		eTokenSpecified := errors.New(errorMsg)
		response.Exception(w, eTokenSpecified, 401)
		return eTokenSpecified
	}

	if !parsedToken.Valid {
		eInvalidToken := errors.New("Token is invalid")
		response.Exception(w, eInvalidToken, 401)
		return eInvalidToken
	}

	newRequest := r.WithContext(context.WithValue(r.Context(), jwtMiddleware.Options.UserProperty, parsedToken))
	*r = *newRequest

	/** check role */
	if role == "" {
		return nil
	}
	tokenRole, _ := ExtractClaim(r, "role")
	userID, _ := ExtractClaim(r, "user_id")
	for k, r := range myRole {
		if k == role {
			for _, c := range r {
				if strings.ToLower(c) == strings.ToLower(tokenRole.(string)) {
					return nil
				}
			}
			break
		}
	}
	e := errors.New("Access is not permitted")
	writeLog.Write(e, "my-role", myRole, "role", role, "token-role", tokenRole, "token-uid", userID, "token", token)
	response.Exception(w, e, 401)
	return e
}
