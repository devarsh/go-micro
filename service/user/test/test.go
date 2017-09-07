package main

import (
	"fmt"
	"github.com/devarsh/micro/service/db"
	"github.com/devarsh/micro/service/user"
	"gopkg.in/mgo.v2"
)

func main() {
	sess, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	x := db.NewDb(sess, "test", user.CollectionName)
	_ = x.DropCollection()
	x.EnsureIndex(user.UniqueKeys)
	x = nil

	Um := user.NewUserMananger(sess, "test")

	err = Um.Create("devarsh", "shah", []string{"admin"})
	fmt.Println("User created", err)

	ok, err := Um.Exists("devarsh")
	fmt.Println("User Exist", ok, err)

	us, err := Um.FindByName("devarsh")
	fmt.Println("user found", us, err)

	err = us.CheckUserCredentials("shah")
	fmt.Println("User password valid", err)

	fmt.Println("User is active", us.Active)
	err = Um.SetActive(us.Username, false)
	fmt.Println("User active status changed", err)

	err = Um.SetPassword(us.Username, "shah1")
	fmt.Println("user password changed", err)

	us, err = Um.FindByName("devarsh")
	fmt.Println("user found", us, err)
	err = us.CheckUserCredentials("test1")
	fmt.Println("User password valid", err)
	fmt.Println("User is active", us.Active)
	fmt.Println("User has claim admin", us.HasClaim("admin"))
	fmt.Println("User has invalid claim", us.HasClaim("invalid"))

	fmt.Println("User is active", us.Active)
	err = Um.SetActive(us.Username, true)
	fmt.Println("User active status changed", err)

	us, err = Um.FindByName("devarsh")
	fmt.Println("user found", us, err)
	fmt.Println("User is active", us.Active)

	err = us.CheckUserCredentials("test1")
	fmt.Println("User password valid", err)

	err = us.CheckUserCredentials("shah1")
	fmt.Println("User password valid", err)

	err = Um.Create("", "shah", []string{"fake"})
	fmt.Println("Empty User created", err)

	err = Um.Create("Devarsh", "shah", []string{"admin"})
	fmt.Println("Duplicate User created", err)

	//Update user not exist
}
