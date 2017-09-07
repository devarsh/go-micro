package token

import (
	"errors"
	"github.com/devarsh/micro/service/user"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Token struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	UserId    bson.ObjectId `bson:"userId"`
	Token     string        `bson:"token"` //unique
	ExpiresAt time.Time     `bson:"expiresAt"`
}

type TokenManager struct {
	duration       time.Duration
	session        *mgo.Session
	dbName         string
	collectionName string
}

var (
	ErrorTokenExpiredOrInactive = errors.New("Token is expired or not yet active")
	ErrorTokenNotFound          = errors.New("Token not found")
	ErrorTokenAlreadyExist      = errors.New("Token already exists")
	ErrorEmptyUserObject        = errors.New("Empty User Object passed or ID is null")
	UniqueKeys                  = []string{"token"}
	CollectionName              = "Token"
)

// expiryDelta determines how earlier a token should be considered
// expired than its actual expiration time. It is used to avoid late
// expirations due to client-server time mismatches.
const ExpiryDelta = 10 * time.Second

func (t *Token) Expired() bool {
	if t.ExpiresAt.IsZero() {
		return false
	}
	return t.ExpiresAt.Add(-ExpiryDelta).Before(time.Now())
}

func NotFoundError(err error) error {
	if err == mgo.ErrNotFound {
		return ErrorTokenNotFound
	}
	return err
}

func NewTokenMananger(duration time.Duration, session *mgo.Session, dbName string) *TokenManager {

	return &TokenManager{duration: duration, session: session, dbName: dbName, collectionName: CollectionName}
}

func (t *TokenManager) Issue(usr *user.User) (string, error) {
	if usr.Id == "" {
		return "", ErrorEmptyUserObject
	}
	tmpToken, err := RandToken()
	if err != nil {
		return "", err
	}
	token := &Token{Id: bson.NewObjectId(),
		UserId:    usr.Id,
		Token:     tmpToken,
		ExpiresAt: time.Now().Add(t.duration),
	}
	err = t.session.DB(t.dbName).C(t.collectionName).Insert(token)
	if err != nil {
		if mgo.IsDup(err) {
			return "", ErrorTokenAlreadyExist
		}
		return "", err
	}
	return tmpToken, nil
}

func (t *TokenManager) Validate(tokenStr string) (*Token, error) {
	token := &Token{}
	if err := t.session.DB(t.dbName).C(t.collectionName).Find(bson.M{"token": tokenStr}).One(token); err != nil {
		return nil, NotFoundError(err)
	}
	if token.Expired() {
		return nil, ErrorTokenExpiredOrInactive
	}
	return token, nil
}

func (t *TokenManager) ForceExpireAll(usr *user.User) error {
	info, err := t.session.DB(t.dbName).C(t.collectionName).UpdateAll(bson.M{"userId": usr.Id, "expiresAt": bson.M{"$gte": time.Now()}}, bson.M{"$set": bson.M{"expiresAt": time.Now().Add(-ExpiryDelta)}})
	if err != nil {
		return err
	}
	if info.Matched > 0 {
		return nil
	}
	return ErrorTokenNotFound
}

func (t *TokenManager) ForeExpireToken(tokenStr string) error {
	err := t.session.DB(t.dbName).C(t.collectionName).Update(bson.M{"token": tokenStr}, bson.M{"$set": bson.M{"expiresAt": time.Now().Add(-ExpiryDelta)}})
	if err != nil {
		return NotFoundError(err)
	}
	return nil
}

func (t *TokenManager) RemoveExpired() error {
	info, err := t.session.DB(t.dbName).C(t.collectionName).RemoveAll(bson.M{"expiresAt": bson.M{"$lte": time.Now()}})
	if err != nil {
		return err
	}
	if info.Matched > 0 {
		return nil
	}
	return ErrorTokenNotFound
}
