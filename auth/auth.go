package auth

import (
	"fmt"
	"github.com/Arnobkumarsaha/new-server/errorhandler"
	"net/http"
	"time"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	creds, err := checkCredentialExist(w, r)
	if err != nil {
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	signedTokenString, err := getSignedTokenString(w, creds.Username, expirationTime)
	if err != nil {
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   signedTokenString,
		Expires: expirationTime,
	})
	errorhandler.Write(w, fmt.Sprintf("%s is signedIn !", creds.Username))
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	tknStr, err := getToken(w, r)
	if err != nil {
		return
	}
	claims, err := getClaimFromTokenString(w, tknStr)
	if err != nil {
		return
	}

	errorhandler.Write(w, fmt.Sprintf("Welcome %s!", claims.Username))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	tknStr, err := getToken(w, r)
	if err != nil {
		return
	}
	claims, err := getClaimFromTokenString(w, tknStr)
	if err != nil {
		return
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 2 Minute of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 2*time.Minute {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	signedTokenString, err := getSignedTokenString(w, claims.Username, expirationTime)
	if err != nil {
		return
	}

	// Set the new token as the userCredentials `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   signedTokenString,
		Expires: expirationTime,
	})
	errorhandler.Write(w, fmt.Sprintf("%s's token has been refresed !", claims.Username))
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	tknStr, err := getToken(w, r)
	if err != nil {
		return
	}

	claims, err := getClaimFromTokenString(w, tknStr)
	if err != nil {
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tknStr,
		Expires: time.Now().Add(1 * time.Microsecond),
	})

	errorhandler.Write(w, fmt.Sprintf("%s is loggedout !", claims.Username))
}
