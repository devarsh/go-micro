package auth

import (
	"github.com/devarsh/micro/service/token"
	"github.com/devarsh/micro/service/token/jwt"
	"github.com/devarsh/micro/service/user"
	"gopkg.in/mgo.v2"
	"time"
)

type AuthService struct {
	userManager  *user.UserManager
	refreshToken *token.TokenManager
	accessToken  *jwt.JwtTokenManager
}

type ApiToken struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthService(session *mgo.Session, refreshTokenTime time.Duration, accessTokenTime time.Duration, tokenKey string) *AuthService {
	um := user.NewUserMananger(session, "test")
	rt := token.NewTokenMananger(refreshTokenTime, session, "test")
	at := jwt.NewJwtTokenManager(accessTokenTime, tokenKey, "localhost:8080", session, "test")
	return &AuthService{userManager: um, refreshToken: rt, accessToken: at}
}

func (as *AuthService) PerformLogin(username, password string) (*ApiToken, error) {
	user, err := as.userManager.FindByName(username)
	if err != nil {
		return nil, err
	}
	err = user.CheckUserCredentials(password)
	if err != nil {
		return nil, err
	}
	rtoken, err := as.refreshToken.Issue(user)
	if err != nil {
		return nil, err
	}
	atoken, err := as.accessToken.Issue(user)
	if err != nil {
		return nil, err
	}
	return &ApiToken{AccessToken: atoken, RefreshToken: rtoken}, nil
}

func (as *AuthService) CreateUser(username, password string, claims []string) error {
	err := as.userManager.Create(username, password, claims)
	if err != nil {
		return err
	}
	return nil
}

func (as *AuthService) ValidateAccessToken(accessToken string) (*jwt.CustomClaims, error) {
	claims, err := as.accessToken.Validate(accessToken)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (as *AuthService) ClearInvalidTokens() error {
	err := as.accessToken.RemoveExpired()
	if err != nil {
		return err
	}
	err = as.refreshToken.RemoveExpired()
	if err != nil {
		return err
	}
	return nil
}

func (as *AuthService) RefreshAccessToken(refreshToken string) (string, error) {
	token, err := as.refreshToken.Validate(refreshToken)
	if err != nil {
		return "", err
	}
	user, err := as.userManager.FindByID(token.Id)
	if err != nil {
		return "", err
	}
	accessTkn, err := as.accessToken.Issue(user)
	if err != nil {
		return "", err
	}
	return accessTkn, nil
}

func (as *AuthService) PerformLogout(username string) error {
	user, err := as.userManager.FindByName(username)
	if err != nil {
		return err
	}
	err = as.refreshToken.ForceExpireAll(user)
	if err != nil {
		return err
	}
	err = as.accessToken.ForceExpireAll(user)
	if err != nil {
		return err
	}
	return nil
}

func (as *AuthService) DeactivateUser(username string) error {
	err := as.userManager.SetActive(username, false)
	if err != nil {
		return err
	}
	err = as.PerformLogout(username)
	if err != nil {
		return err
	}
	return nil
}
