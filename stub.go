package main

import (
	"fmt"
	//"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Demo struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	UserId    string
	Active    bool
	LoginTime time.Time
}

func newDemo(userid string) *Demo {
	tm := time.Now()
	return &Demo{Id: bson.NewObjectId(), UserId: userid, Active: true, LoginTime: tm}
}

func main() {
	/*
		data, err := bcrypt.GenerateFromPassword([]byte("devarsh"), 2)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("%s\n", data)
		err = bcrypt.CompareHashAndPassword([]byte("$2a$10$S.fFnHQ.K0GEpJTkxCBvW.8LUfij65Ttx0..WqN7/ljsTNPTz.IX2"), []byte("@#$$!#@$@QRWDASFCDSX Z%CSA"))
		if err != nil {
			fmt.Print(err)
		}*/

	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	//demo := &Demo{}
	/*err = session.DB("test").C("demo").Update(bson.M{"userid": "Harsh"}, bson.M{"logintime": time.Now()})
	if err != nil {
		fmt.Print(err)
	} */
	//Insert
	/*dm := newDemo("devarsh")
	session.DB("test").C("demo").Insert(dm)
	time.Sleep(time.Second * 10)
	dm = newDemo("devarsh")
	session.DB("test").C("demo").Insert(dm)
	time.Sleep(time.Second * 20)
	dm = newDemo("Harsh")
	session.DB("test").C("demo").Insert(dm)
	time.Sleep(time.Second * 30)
	dm = newDemo("Harsh")
	session.DB("test").C("demo").Insert(dm)
	*/
	//Insert Session Close
	/*
		dm := newDemo("DEva")
		err = session.DB("test").C("demo").Insert(dm)
		if err != nil {
			panic(err)
		}
		err = session.DB("test").C("demo").Insert(dm)
		if err != nil {
			fmt.Println(mgo.IsDup(err))
		}
		fmt.Println("This is the end.")
		//Bulk Update
		/*
			loc, err := time.LoadLocation("Asia/Kolkata")
			if err != nil {
				panic(err)
			}
			fmt.Println(loc)
			tm := time.Date(2017, time.August, 16, 15, 38, 45, 0, loc)
			fmt.Println(tm)
			info, err := session.DB("test").C("demo").UpdateAll(bson.M{"logintime": bson.M{"$gt": tm}}, bson.M{"$set": bson.M{"active": true}})
			if err != nil {
				panic(err)
			}
			fmt.Println(info.Matched, info.Updated, info.Removed, info.UpsertedId)
	*/
	//Single Update
	/*err = session.DB("test").C("demo").Update(bson.M{"userid": bson.RegEx{"DEVARSHG", "i"}}, bson.M{"$set": bson.M{"logintime": time.Now()}})
	if err != nil {
		if err == mgo.ErrNotFound {
			fmt.Println("Oops no data found error")
		} else {
			panic(err)
		}
	}
	fmt.Print("updated")
	*/
	//Find With Count
	/*
		cnt, err := session.DB("test").C("demo").Find(bson.M{"userid": bson.RegEx{"DEVARSHG", "i"}}).Count()
		if err != nil {
			panic(err)
		}

		fmt.Println(cnt)
	*/
	//Find One
	/*dm := &Demo{}
	err = session.DB("test").C("demo").Find(bson.M{"userid": bson.RegEx{"DEVARSHG", "i"}}).One(dm)
	if err == mgo.ErrNotFound {
		fmt.Print("Oops no data")
	}*/

}
