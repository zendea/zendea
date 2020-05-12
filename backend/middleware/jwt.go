package middleware

import (
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"zendea/form"
	"zendea/model"
	"zendea/service"
	"zendea/util/log"
)

//login type
var (
	LoginStandard = 1
	LoginOAuth    = 2
)

type LoginDto struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Code     string `form:"code" json:"code"`
}

// LoginOAuthDto - oauth login
type LoginOAuthDto struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}

//todo : 用单独的claims model去掉user model
func JwtAuth(LoginType int) *jwt.GinJWTMiddleware {
	jwtMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "Jwt",
		// SigningAlgorithm: "RS256",
		// PubKeyFile:       "keys/jwt_private_key.pem",
		// PrivKeyFile:      "keys/jwt_public_key.pem",
		Key:              []byte(viper.GetString("jwt.key")),
		Timeout:          time.Hour * 24,
		MaxRefresh:       time.Hour * 24 * 90,
		IdentityKey:      viper.GetString("jwt.identity_key"),
		LoginResponse:    LoginResponse,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(model.UserClaims); ok {
				return jwt.MapClaims{
					"id":    v.ID,
					"name":  v.Name,
					"uid":   v.ID,
					"uname": v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return model.UserClaims{
				Name: claims["name"].(string),
				ID:   int64(claims["id"].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			if LoginType == LoginOAuth { //OAuth
				return AuthenticatorOAuth(c)
			}
			return Authenticator(c)
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(model.UserClaims); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(200, gin.H{
				"code":    200,
				"success": false,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Error(err.Error())
	}
	return jwtMiddleware
}

func LoginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": map[string]interface{}{
			"token":  token,
			"expire": expire,
		},
		"success": true,
		"message": "success",
	})
}
func Authenticator(c *gin.Context) (interface{}, error) {
	var loginDto LoginDto
	if err := form.Bind(c, &loginDto); err != nil {
		return "", err
	}

	log.Info("loginDto.Username: %s", loginDto.Username)

	ok, err, u := service.UserService.VerifyAndReturnUserInfo(loginDto.Username, loginDto.Password) // Standard login
	if ok {

		return model.UserClaims{
			ID:   u.ID,
			Name: u.Username.String,
		}, nil
	}
	return nil, err
}

func AuthenticatorOAuth(c *gin.Context) (interface{}, error) {
	provider := c.Param("provider")

	var oauthDto LoginOAuthDto
	if err := form.Bind(c, &oauthDto); err != nil {
		return "", err
	}

	account, err := service.LoginSourceService.GetOrCreate(provider, oauthDto.Code, oauthDto.State)
	if err != nil {
		
		return nil, err
	}

	u, err := service.UserService.SignInByLoginSource(account)
	if err == nil {
		return model.UserClaims{
			ID:   u.ID,
			Name: u.Username.String,
		}, nil
	}

	log.Info("oauthDto.Code: %s", oauthDto.Code)
	log.Info("oauthDto.State: %s", oauthDto.State)
	return nil, err
}
