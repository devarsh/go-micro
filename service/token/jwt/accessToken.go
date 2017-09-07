package jwt

import (
	"errors"
	"fmt"
	Rtoken "github.com/devarsh/micro/service/token"
	"github.com/devarsh/micro/service/user"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type AccessToken struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Jwt       string        `bson:"jwt"` //unique
	ExpiresAt time.Time     `bson:"expiresAt"`
	UserId    bson.ObjectId `bson:"userId"`
	Blacklist bool          `bson:"blacklist"`
}

type CustomClaims struct {
	UserId bson.ObjectId
	jwt.StandardClaims
}

type JwtTokenManager struct {
	tokenDuration  time.Duration // in minutes
	privateKey     []byte
	issuer         string
	session        *mgo.Session
	collectionName string
	dbName         string
}

var (
	ErrorInvalidOrMalformedToken   = errors.New("Invalid or Malformed Token")
	ErrorTokenExpiredOrInactive    = errors.New("Token is expired or not yet active")
	ErrorInvalidTokenSigningMethod = errors.New("Invalid token signing method")
	ErrorTokenNotFound             = errors.New("Token not found")
	ErrorTokenAlreadyExist         = errors.New("Token already exists")
	ErrorEmptyUserObject           = errors.New("Empty User Object passed or ID is null")
	UniqueKeys                     = []string{"jwt"}
	CollectionName                 = "AccessToken"
)

func NotFoundError(err error) error {
	if err == mgo.ErrNotFound {
		return ErrorTokenNotFound
	}
	return err
}

func NewJwtTokenManager(tokenDuration time.Duration, privateKey string, issuer string, sess *mgo.Session, dbName string) *JwtTokenManager {
	jwtM := &JwtTokenManager{}
	jwtM.session = sess
	jwtM.collectionName = CollectionName
	jwtM.dbName = dbName
	jwtM.tokenDuration = tokenDuration
	jwtM.privateKey = []byte(privateKey)
	jwtM.issuer = issuer
	return jwtM
}

func (j *JwtTokenManager) Issue(usr *user.User) (string, error) {
	if usr.Id == "" {
		return "", ErrorEmptyUserObject
	}
	nowTime := time.Now()
	expiresAt := nowTime.Add(j.tokenDuration)
	uid, err := Rtoken.RandToken()
	if err != nil {
		return "", err
	}
	claims := CustomClaims{usr.Id,
		jwt.StandardClaims{
			IssuedAt:  nowTime.Unix(),
			ExpiresAt: expiresAt.Unix(),
			Id:        uid,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}
	accessToken := &AccessToken{Id: bson.NewObjectId(), Jwt: ss, UserId: usr.Id, ExpiresAt: expiresAt, Blacklist: false}
	err = j.session.DB(j.dbName).C(j.collectionName).Insert(accessToken)
	if err != nil {
		if mgo.IsDup(err) {
			return "", ErrorTokenAlreadyExist
		}
		return "", err
	}
	return ss, nil
}

func (j *JwtTokenManager) isJwtValid(ss string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(ss, &CustomClaims{}, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrorInvalidTokenSigningMethod
		}
		return j.privateKey, nil
	})
	fmt.Println("---------", token.Valid, "------------")
	if token.Valid {

		if claims, ok := token.Claims.(*CustomClaims); ok {
			return claims, nil
		}
		return nil, ErrorInvalidOrMalformedToken
	}
	ve, ok := err.(*jwt.ValidationError)
	if ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, ErrorInvalidOrMalformedToken
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, ErrorTokenExpiredOrInactive
		}
	}
	return nil, err
}

func (j *JwtTokenManager) Validate(ss string) (*CustomClaims, error) {
	token, err := j.isJwtValid(ss)
	if err != nil {
		return nil, err
	}
	accessToken := &AccessToken{}
	err = j.session.DB(j.dbName).C(j.collectionName).Find(bson.M{"jwt": ss}).One(accessToken)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrorTokenExpiredOrInactive
		}
		return nil, err
	}
	if accessToken.Blacklist == true {
		return nil, ErrorTokenExpiredOrInactive
	}
	return token, nil
}

func (j *JwtTokenManager) ForceExpireAll(usr *user.User) error {
	info, err := j.session.DB(j.dbName).C(j.collectionName).UpdateAll(bson.M{"userId": usr.Id, "blacklist": false, "expiresAt": bson.M{"$gte": time.Now()}}, bson.M{"$set": bson.M{"blacklist": true}})
	if err != nil {
		return err
	}
	if info.Matched > 0 {
		return nil
	}
	return ErrorTokenNotFound
}

func (j *JwtTokenManager) ForeExpireToken(ss string) error {
	err := j.session.DB(j.dbName).C(j.collectionName).Update(bson.M{"jwt": ss}, bson.M{"$set": bson.M{"blacklist": true}})
	if err != nil {
		if err == mgo.ErrNotFound {
			return ErrorTokenNotFound
		}
		return err
	}
	return nil
}

func (j *JwtTokenManager) RemoveExpired() error {
	info, err := j.session.DB(j.dbName).C(j.collectionName).RemoveAll(bson.M{"$or": []bson.M{bson.M{"expiresAt": bson.M{"$lte": time.Now()}}, bson.M{"blacklist": true}}})
	if err != nil {
		return err
	}
	if info.Matched > 0 {
		return nil
	}
	return ErrorTokenNotFound
}
