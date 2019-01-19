package api

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/openshift/osin"
	"github.com/pmoncadaisla/istio-auth-sample/authn/pkg/oauth2/dto"
	"github.com/pmoncadaisla/istio-auth-sample/authn/pkg/oauth2/services"
)

type api struct {
	server      *osin.Server
	authService services.AuthenticatorValidator
}

// RestHandler represent api methods or behavior. What a oauth2 api do.
type RestHandler interface {
	SetupRouting(g *gin.RouterGroup)
	MakePong() func(*gin.Context)
}

var once sync.Once

// Instance represent a single api instance.
var Instance *api

// New is an api factory method (sigleton)
func New(osinStorage osin.Storage, accessTokenExpiration int32, authServ services.AuthenticatorValidator) RestHandler {
	once.Do(func() {
		Instance = new(api)
		sconfig := osin.NewServerConfig()
		sconfig.AccessExpiration = accessTokenExpiration
		sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
		sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
			osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}
		sconfig.AllowGetAccessRequest = true
		sconfig.AllowClientSecretInParams = true
		//sconfig.RequirePKCEForPublicClients = true // with a clientID and secret it's enough
		Instance.server = osin.NewServer(sconfig, osinStorage)
		Instance.server.AccessTokenGen = authServ
		Instance.authService = authServ
	})

	return Instance
}

func (d *api) SetupRouting(g *gin.RouterGroup) {
	g.PUT("/ping", d.MakePong())
	g.POST("/token", d.MakeAccessToken())
	g.GET("/jwks", d.RetrievePubKeys())
}

func (o *api) MakePong() func(*gin.Context) {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(200, gin.H{
			"response": "pong",
		})

		return
	}
}

func (o *api) RetrievePubKeys() func(*gin.Context) {
	return func(c *gin.Context) {
		pubKey := o.authService.RetrieveJWK()
		keys := dto.KeysNew(pubKey)

		c.AbortWithStatusJSON(200, keys)
		return
	}
}

func (o *api) MakeAccessToken() func(c *gin.Context) {
	return func(c *gin.Context) {

		//requestID := c.Keys["uid"].(string)
		resp := o.server.NewResponse()
		defer resp.Close()
		c.Request.ParseForm()

		if ar := o.server.HandleAccessRequest(resp, c.Request); ar != nil {
			switch ar.Type {
			case osin.PASSWORD:
				if o.authService.ExistUser(ar.Username) {
					o.authService.LoadUserData(ar, false)
					ar.Authorized = true

					var uData dto.CustomUserData
					uData.Name = ar.Username

					var accessData osin.AccessData
					accessData.UserData = uData
					ar.AccessData = &accessData
				} else {
					// clientID and secret are validated in HandleAccessRequest
					//ar.Authorized = true // Not allowed yet
					ar.Authorized = false
				}
				break
			default:
				c.AbortWithStatusJSON(500, gin.H{
					"error": "unexpected grant type",
				})
				return
			}

			logrus.Info(ar.AccessData.UserData)
			o.server.FinishAccessRequest(resp, c.Request, ar)
			osin.OutputJSON(resp, c.Writer, c.Request)
			return
		}

		if resp.IsError && resp.InternalError != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error": "unexpected error",
			})
			return
		}

		return

	}
}
