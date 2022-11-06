package auth

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

// Decode the credentials from request.Body & check whether it is ok or not.
func checkCredentialExist(w http.ResponseWriter, r *http.Request) (*Credentials, error) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP errorhandler
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	fmt.Println("Credentials passed = ", creds)

	// Get the expected password from our in memory map
	expectedPassword, ok := userCredentials[creds.Username]

	// If a password exists for the given user AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, err
	}
	return &creds, nil
}

// Read the tokenString from request.Cookie, with some error handling
func getToken(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return "", err
		}
		// For any other type of errorhandler, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return "", err
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	fmt.Printf("Token string : %v %v", tknStr, c)
	return tknStr, nil
}

func getClaimFromTokenString(w http.ResponseWriter, tknStr string) (*Claims, error) {
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well.
	// It returns error if the token is invalid, expired or signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(*jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	fmt.Printf("Claim : %v", claims)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return nil, err
		}
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, err
	}
	return claims, nil
}

func getSignedTokenString(w http.ResponseWriter, username string, expirationTime time.Time) (string, error) {
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an errorhandler in creating the JWT return an internal server errorhandler
		w.WriteHeader(http.StatusInternalServerError)
		return "", err
	}
	fmt.Println("tokenString = ", tokenString)
	return tokenString, nil
}
