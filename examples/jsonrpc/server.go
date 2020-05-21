package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
	"github.com/gorilla/sessions"

	"github.com/nobonobo/wecty/examples/jsonrpc/frontend/models"
)

const SESSION_KEY = "hodghu37dbegdyu37"

var (
	users = map[string]*models.User{
		"nobo@sample.com": {EMail: "nobo@sample.com", Name: "nobonobo", Password: "piyopiyo"},
	}
	store = sessions.NewCookieStore([]byte(SESSION_KEY))
)

// Service ...
type Service struct{}

// User ...
func (s *Service) User(r *http.Request, none *struct{}, rep *models.User) error {
	println("call Service.User")
	session, _ := store.Get(r, "session")
	email, ok := session.Values["email"].(string)
	if !ok || len(email) == 0 {
		return fmt.Errorf("user not found")
	}
	if u, ok := users[email]; ok {
		*rep = models.User{
			EMail: u.EMail,
			Name:  u.Name,
		}
		return nil
	}
	return fmt.Errorf("user not found")
}

func main() {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(&Service{}, "")
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1024); err != nil {
			log.Print(err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		session, _ := store.Get(r, "session")
		email := r.FormValue("email")
		log.Println("login:", email)
		if u, ok := users[email]; ok {
			if u.Password == r.FormValue("password") {
				session.Values["email"] = email
				session.Values["authenticated"] = true
				session.Save(r, w)
				return
			}
		}
		err := fmt.Errorf("mismatch email or password")
		log.Print(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
	})
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		// Revoke users authentication
		delete(session.Values, "email")
		session.Values["authenticated"] = false
		session.Save(r, w)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		log.Print(r)
		s.ServeHTTP(w, r)
		/*
			s.ServeRequest(
				jsonrpc.NewServerCodec(
					&struct {
						io.ReadCloser
						io.Writer
					}{
						ioutil.NopCloser(r.Body),
						w,
					},
				),
			)
		*/
	})
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}
