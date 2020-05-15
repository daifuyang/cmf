package router

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gincmf/cmf/controller"
	"golang.org/x/oauth2"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"net/http"
	"time"
)

var (
	AppPort *string
	authServerURL = "http://localhost:"
	globalToken *oauth2.Token
	rc controller.RestControllerStruct
	Srv *server.Server
)

func RegisterOauthRouter(e *gin.Engine) {

	fmt.Println("AppPort",*AppPort)

	authServerURL += *AppPort

	h := md5.New()
	h.Write([]byte("gincmf"))
	md5str := hex.EncodeToString(h.Sum(nil))

	config := oauth2.Config{
		ClientID:     "1",
		ClientSecret: md5str,
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

	accessToken := generates.NewJWTAccessGenerate([]byte("gincmf"), jwt.SigningMethodHS512)
	manager.MapAccessGenerate(accessToken)
	clientStore := store.NewClientStore()
	clientStore.Set(config.ClientID, &models.Client{
		ID:     config.ClientID,
		Secret: config.ClientSecret,
		Domain: authServerURL,
	})

	manager.MapClientStorage(clientStore)
	Srv = server.NewServer(server.NewConfig(), manager)
	Srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "admin" && password == "123456" {
			 userID = "admin"
		}
		return userID,nil
	})

	e.GET("api/oauth/token",func(c *gin.Context){

		username := c.Query("username")
		password := c.Query("password")

		fn := Srv.PasswordAuthorizationHandler
		userID,err := fn(username,password)

		if userID == "" {
			rc.Error(c,"userid不能为空！")
			return
		}

		req := &server.AuthorizeRequest{
			RedirectURI:  config.RedirectURL,
			ResponseType: "code",
			ClientID:     config.ClientID,
			State:        "jwt",
			Scope:        "all",
			UserID:         userID,
			AccessTokenExp: 5 * time.Hour,
			Request:      c.Request,
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
		c.JSON(http.StatusOK,token)
		})

	e.POST("/token",func(c *gin.Context) {
		fmt.Println("token request")
		err := Srv.HandleTokenRequest(c.Writer, c.Request)
		if err != nil {
			panic(err.Error())
		}
	})
}