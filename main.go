package main

import (
	"github.com/devarsh/micro/service/auth"
	"gopkg.in/mgo.v2"
	"net/http"
)

func main() {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	authServ := auth.NewAuthService(session)

	http.ListenAndServe(":8081", nil)
}

func AddUser(sess *mgo.Session) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
