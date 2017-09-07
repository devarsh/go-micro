package main

import (
	"fmt"
	"github.com/devarsh/micro/service/db"
	"github.com/devarsh/micro/service/token"
	"github.com/devarsh/micro/service/user"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	x = nil

	usr := &user.User{
		Id:       bson.NewObjectId(),
		Username: "devarsh",
		Password: []byte("adminadmin"),
		Claims:   []string{"admin"},
	}

	jm := token.NewTokenMananger(time.Second*time.Duration(15), sess, "test")

	if jm == nil {
		panic("Jwt manager is nil")
	}
	token, err := jm.Issue(usr)
	fmt.Println("Token1 generated : ", token, err)
	valid, err := jm.Validate(token)
	fmt.Println("Token1 valid : ", valid, err)
	err = jm.ForeExpireToken(token)
	fmt.Println("Token1 force expired, if no err its successful, err == : ", err)

	valid, err = jm.Validate(token)
	fmt.Println("Token1 revalidate afer expiry : ", valid, err)

	token2, err := jm.Issue(usr)
	fmt.Println("Token2 generated : ", token2, err)

	fmt.Println("Will sleep for 10 second : ")
	time.Sleep(time.Second * time.Duration(11))

	valid, err = jm.Validate(token2)
	fmt.Println("Token 2 should expires after expiry time : ", valid, err)

	token3, err := jm.Issue(usr)
	fmt.Println("Token3 generated : ", token3, err)
	token4, err := jm.Issue(usr)
	fmt.Println("Token4 generated : ", token4, err)

	err = jm.ForceExpireAll(usr)
	fmt.Println("All tokens for user expired : ", err)

	valid, err = jm.Validate(token3)
	fmt.Println("Token3 valid : ", valid, err)
	valid, err = jm.Validate(token4)
	fmt.Println("Token3 valid : ", valid, err)

	token5, err := jm.Issue(usr)
	fmt.Println("Token5 generated : ", token5, err)

	err = jm.RemoveExpired()
	fmt.Println("Delete expired tokens from the db : ", err)

	valid, err = jm.Validate(token3)
	fmt.Println("Token3 valid after delete : ", valid, err)
	valid, err = jm.Validate(token4)
	fmt.Println("Token3 valid after delete : ", valid, err)

}
