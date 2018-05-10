package controller

import (
	"example-go-web/util"
	"net/http"
	"net/url"
	"time"
)

type Session struct {
	ID     string
	UserID string
	Expiry time.Time
}

const (
	// Keep users logged in for 3 days
	sessionLength     = 24 * 3 * time.Hour
	sessionCookieName = "GoesSession"
	sessionIDLength   = 20
)

/*
	mengembalikan waktu sebelum Session Expiry
*/
func (session *Session) Expired() bool {
	return session.Expiry.Before(time.Now())
}

/*
	untuk mendaftarkan session baru bagi user yang login
*/
func NewSession(w http.ResponseWriter) *Session {
	// berapa lama sebuah session di-simpan
	expiry := time.Now().Add(sessionLength)

	// inisialisasi session untuk user yg login/daftar
	session := &Session{
		ID:     util.GenerateID("sess", sessionIDLength),
		Expiry: expiry,
	}

	// mendaftarkan session pada cookie untuk di-ingat
	cookie := http.Cookie{
		Name:    sessionCookieName,
		Value:   session.ID,
		Expires: session.Expiry,
	}

	// adding session pada cookie
	http.SetCookie(w, &cookie)
	return session
}

/*
	find session dari cookie.value
*/
func RequestSession(r *http.Request) *Session {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil
	}

	session, err := globalSessionStore.Find(cookie.Value)
	if err != nil {
		panic(err)
	}

	if session == nil {
		return nil
	}

	if session.Expired() {
		globalSessionStore.Delete(session)
		return nil
	}
	return session
}

func RequireLogin(w http.ResponseWriter, r *http.Request) {
	// Let the request pass if we've got a user
	if RequestUser(r) != nil {
		return
	}

	query := url.Values{}
	query.Add("next", url.QueryEscape(r.URL.String()))

	http.Redirect(w, r, "/login?"+query.Encode(), http.StatusFound)
}

func RequestUser(r *http.Request) *User {
	session := RequestSession(r)
	if session == nil || session.UserID == "" {
		return nil
	}

	user, err := GlobalUserStore.Find(session.UserID)
	if err != nil {
		panic(err)
	}
	return user
}

func FindOrCreateSession(w http.ResponseWriter, r *http.Request) *Session {
	session := RequestSession(r)
	if session == nil {
		session = NewSession(w)
	}

	return session
}
