package util

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/magicalsoup/reelgo/.gen/reelgo/public/model"
)

func GetBearerToken(cookies []*http.Cookie) (string, error) {
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == BEARER_TOKEN_COOKIE_NAME {
			return cookies[i].Value, nil
		}
	}
	return "", errors.New(BEARER_TOKEN_COOKIE_NAME + " cookie not found ")
}

const BEARER_TOKEN_COOKIE_NAME = "user-bearer-token"
const USER_ID_COOKIE_NAME = "user-id"

func SetAuthCookies(w http.ResponseWriter, token *model.Tokens) {
	id_cookie := http.Cookie{
		Name:     USER_ID_COOKIE_NAME,
		Value:    strconv.Itoa(int(token.UID)),
		Path:     "/",
		MaxAge:   int(token.ExpiryTime),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	bearer_cookie := http.Cookie{
		Name:     BEARER_TOKEN_COOKIE_NAME,
		Value:    token.BearerToken,
		Path:     "/",
		MaxAge:   int(token.ExpiryTime),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &id_cookie)
	http.SetCookie(w, &bearer_cookie)
}