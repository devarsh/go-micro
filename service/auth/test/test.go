package main

import (
	"fmt"
	"github.com/devarsh/micro/service/auth"
	"github.com/devarsh/micro/service/db"
	"github.com/devarsh/micro/service/token"
	"github.com/devarsh/micro/service/token/jwt"
	"github.com/devarsh/micro/service/user"
	"gopkg.in/mgo.v2"
	"time"
)

func main() {
	sess, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	x := db.NewDb(sess, "test", token.CollectionName)
	_ = x.DropCollection()
	x.EnsureIndex(token.UniqueKeys)
	y := db.NewDb(sess, "test", user.CollectionName)
	_ = y.DropCollection()
	y.EnsureIndex(user.UniqueKeys)
	z := db.NewDb(sess, "test", jwt.CollectionName)
	_ = z.DropCollection()
	z.EnsureIndex(jwt.UniqueKeys)
	x = nil
	y = nil
	z = nil

	am := auth.NewAuthService(sess)

	err = am.CreateUser("Devarsh", "Devarsh123", []string{"admin", "read"})
	fmt.Println("User created successfully", err)

	tokens, err := am.PerformLogin("devarsh", "Devarsh123")
	fmt.Println("User successfully logged", tokens, err)

	time.Sleep(100)

	rtoken, err := am.RefreshAccessToken(tokens.AccessToken)
	fmt.Println("User successfully logged In", rtoken, err)

}
