package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

var (
	ErrorUserNotFound            = errors.New("User Not Found")
	ErrorUserAlreadyExist        = errors.New("User Already Exists")
	ErrorUserInactvie            = errors.New("User is Inactive")
	ErrorUserPasswordInvalid     = errors.New("Username or password is Invalid")
	ErrorEmptyUsernameOrPassword = errors.New("Empty Username or Password")
	UniqueKeys                   = []string{"username"}
	CollectionName               = "Users"
)

type User struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Username string        `bson:"username"` //unique
	Password []byte        `bson:"password"`
	Active   bool          `bson:"active"`
	Claims   []string      `bson:"claims"`
}

func (user *User) HasClaim(claim string) bool {
	for _, oneClaim := range user.Claims {
		if oneClaim == claim {
			return true
		}
	}
	return false
}

func (user *User) CheckUserCredentials(password string) error {
	if user.Active == false {
		return ErrorUserInactvie
	}
	hashPassword := user.Password
	err := bcrypt.CompareHashAndPassword(hashPassword, []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrorUserPasswordInvalid
		}
		return err
	}
	return nil
}

type UserManager struct {
	session        *mgo.Session
	dbName         string
	collectionName string
}

func NotFoundError(err error) error {
	if err == mgo.ErrNotFound {
		return ErrorUserNotFound
	}
	return err
}

func NewUserMananger(session *mgo.Session, dbName string) *UserManager {
	userM := &UserManager{}
	userM.session = session
	userM.collectionName = CollectionName
	userM.dbName = dbName
	return userM
}

func (u *UserManager) Create(username, password string, claims []string) error {
	if username == "" || password == "" {
		return ErrorEmptyUsernameOrPassword
	}
	username = strings.ToLower(username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	if err != nil {
		return err
	}
	user := &User{Id: bson.NewObjectId(),
		Username: username,
		Password: hashedPassword,
		Claims:   claims,
		Active:   true,
	}
	err = u.session.DB(u.dbName).C(u.collectionName).Insert(user)
	if err != nil {
		if mgo.IsDup(err) {
			return ErrorUserAlreadyExist
		}
		return err
	}
	return nil
}

func (u *UserManager) Exists(username string) (bool, error) {
	username = strings.ToLower(username)
	cnt, err := u.session.DB(u.dbName).C(u.collectionName).Find(bson.M{"username": username}).Count()
	if err != nil {
		return false, err
	}
	if cnt > 0 {
		return true, nil
	}
	return false, nil
}

func (u *UserManager) FindByName(username string) (*User, error) {
	username = strings.ToLower(username)
	user := User{}
	err := u.session.DB(u.dbName).C(u.collectionName).Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return nil, NotFoundError(err)
	}
	return &user, nil
}

func (u *UserManager) FindByID(id bson.ObjectId) (*User, error) {
	user := User{}
	err := u.session.DB(u.dbName).C(u.collectionName).FindId(id).One(&user)
	if err != nil {
		return nil, NotFoundError(err)
	}
	return &user, nil
}

func (u *UserManager) SetActive(username string, state bool) error {
	username = strings.ToLower(username)
	err := u.session.DB(u.dbName).C(u.collectionName).Update(bson.M{"username": username}, bson.M{"$set": bson.M{"active": state}})
	if err != nil {
		return NotFoundError(err)
	}
	return nil
}

func (u *UserManager) SetPassword(username, newPassword string) error {
	username = strings.ToLower(username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), -1)
	if err != nil {
		return err
	}
	err = u.session.DB(u.dbName).C(u.collectionName).Update(bson.M{"username": username}, bson.M{"$set": bson.M{"password": hashedPassword}})
	if err != nil {
		return NotFoundError(err)
	}
	return nil
}
