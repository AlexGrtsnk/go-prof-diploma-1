package cookies

import (
	"errors"
	"net/http"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func SetCookieHandler(w http.ResponseWriter, r *http.Request, token string) (cks *http.Cookie) {
	// Initialize the cookie as normal.
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    token,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	// Write the cookie. If there is an error (due to an encoding failure or it
	// being too long) then log the error and send a 500 Internal Server Error
	// response.
	http.SetCookie(w, &cookie)
	/*
		err = Write(w, cookie)
		if err != nil {
			return err
		}
	*/
	return &cookie
	//w.Write([]byte("cookie set!"))
}

func GetCookieHandler(w http.ResponseWriter, r *http.Request) (token string, err error) {
	// Use the Read() function to retrieve the cookie value, additionally
	// checking for the ErrInvalidValue error and handling it as necessary.
	cookie, err := r.Cookie("exampleCookie")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
