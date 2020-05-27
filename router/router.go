package router

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gincmf/cmf/controller"
	"github.com/gincmf/cmf/util"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"net/http"
	"strconv"
	"time"
)

var (
	AppPort       *string
	AuthCode      *string
	Db            *gorm.DB
	authServerURL = "http://localhost:"
	globalToken   *oauth2.Token
	rc            controller.RestController
	Srv           *server.Server
)

type User struct {
	Id           int
	UserType     int
	Gender       int
	Birthday     int
	UserLogin    string `gorm:"type:varchar(60);not null"`
	UserPass     string `gorm:"type:varchar(64);not null"`
	UserNickname string
	Avatar       string
	Signature    string
	Mobile       string
}

func RegisterOauthRouter(e *gin.Engine) {

	authServerURL += *AppPort

	clientSecret := util.GetMd5(*AuthCode)

	fmt.Println("clientSecret",clientSecret)

	config := oauth2.Config{
		ClientID:     "1",
		ClientSecret: clientSecret,
		Scopes:       []string{"all"},
		RedirectURL:  authServerURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/authorize",
			TokenURL: authServerURL + "/token",
		},
	}

	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	accessToken := generates.NewJWTAccessGenerate([]byte("gincmf"+*AuthCode), jwt.SigningMethodHS512)
	manager.MapAccessGenerate(accessToken)
	clientStore := store.NewClientStore()
	clientStore.Set(config.ClientID, &models.Client{
		ID: config.ClientID,
		Secret: config.ClientSecret,
		Domain: authServerURL,
	})

	manager.MapClientStorage(clientStore)
	Srv = server.NewServer(server.NewConfig(), manager)
	Srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		u := &User{}
		userResult := Db.First(u, "user_login = ?", username) // 查询
		userID = ""
		if !userResult.RecordNotFound() {
			//验证密码
			if util.GetMd5(password) == u.UserPass {
				userID = u.UserLogin
			}
		}
		return userID, nil
	})

	e.POST("api/oauth/token", func(c *gin.Context) {

		username := c.PostForm("username")
		password := c.PostForm("password")
		tokenExp := c.DefaultPostForm("expire","2")

		exp,err := strconv.Atoi(tokenExp)

		if err != nil {
			fmt.Println("err",err.Error())
			rc.Error(c, "失效时间应该是整数，单位为小时！")
			return
		}

		fmt.Println("username",username)
		fmt.Println("password",password)

		fn := Srv.PasswordAuthorizationHandler
		userID, err := fn(username, password)

		if userID == "" {
			rc.Error(c, "账号密码不正确！")
			return
		}

		req := &server.AuthorizeRequest{
			RedirectURI:    config.RedirectURL,
			ResponseType:   "code",
			ClientID:       config.ClientID,
			State:          "jwt",
			Scope:          "all",
			UserID:         userID,
			AccessTokenExp: time.Duration(exp) * time.Hour,
			Request:        c.Request,
		}

		ti, err := Srv.GetAuthorizeToken(req)
		if err != nil {
			fmt.Println(err.Error())
		}

		code := ti.GetCode()

		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			//panic(err.Error())
		}

		globalToken = token
		c.JSON(http.StatusOK, token)
	})

	e.POST("api/oauth/refresh", func(c *gin.Context) {
		if globalToken == nil {
			rc.Error(c, "非法访问！")
			return
		}

		username := c.PostForm("username")
		password := c.PostForm("password")

		fn := Srv.PasswordAuthorizationHandler
		userID, err := fn(username, password)

		if userID == "" {
			rc.Error(c, "账号密码不正确！")
			return
		}

		globalToken.Expiry = time.Now()
		token, err := config.TokenSource(context.Background(), globalToken).Token()
		if err != nil {
			rc.Error(c, err.Error())
		}

		globalToken = token
		e := json.NewEncoder(c.Writer)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	e.POST("/token", func(c *gin.Context) {
		fmt.Println("token request")
		err := Srv.HandleTokenRequest(c.Writer, c.Request)
		if err != nil {
			panic(err.Error())
		}
	})
}
