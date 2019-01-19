package oauth2

import (
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pmoncadaisla/istio-auth-sample/authn/pkg/configuration"
	"github.com/pmoncadaisla/istio-auth-sample/authn/pkg/oauth2/api"
	"github.com/pmoncadaisla/istio-auth-sample/authn/pkg/oauth2/services"
	osinStorage "github.com/pmoncadaisla/istio-auth-sample/authn/pkg/oauth2/storage"
	"go.opencensus.io/plugin/ochttp"
)

type Oauth2 struct {
	ContextName   string
	ServerAddress string
}

// HttpServer define Oauth2 server behavior
type HttpServer interface {
	Run()
	createRouter() *gin.Engine
}

var once sync.Once

// Instance represents a single Oauth2 instance
var Instance *Oauth2

// New method is a Oauth2 factory method. If an instance already exist then this instance will be returned (sigleton).
func New(contextName, serverAddress string) HttpServer {
	once.Do(func() {
		Instance = new(Oauth2)
		Instance.ContextName = contextName
		Instance.ServerAddress = serverAddress
	})

	return Instance
}

// Run method is the main Oauth2 method in order to launch a Oauth2
func (d *Oauth2) Run() {
	router := d.createRouter()

	router.Run(d.ServerAddress)
}

func (d *Oauth2) createRouter() (r *gin.Engine) {
	r = customEngine("/"+d.ContextName+"/v1/healthz", "/"+d.ContextName+"/v1/health", "/")
	r.Use(addUniqueRequestID)
	osinStorageInstance := osinStorage.NewTestStorage()

	oauthApi := api.New(osinStorageInstance, configuration.Instance.AccessTokenExpiration, services.New())

	r.Use(addUniqueRequestID)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Headers", "authorization, Content-Type, Authorization")
		c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, HEAD, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}

		if c.Request.RequestURI == "/"+d.ContextName+"/v1/healthz" || c.Request.RequestURI == "/" {
			c.AbortWithStatusJSON(200, gin.H{
				"status": "ok",
			})
		}
	})

	oauthApi.SetupRouting(r.Group("/" + d.ContextName + "/v1/oauth2"))

	return
}

func customEngine(excludePaths ...string) *gin.Engine {
	engine := gin.New()

	if configuration.Instance.TracingEnable {
		engine.Use(gin.WrapH(&ochttp.Handler{}))
	}
	engine.Use(gin.LoggerWithWriter(gin.DefaultWriter, excludePaths...), gin.Recovery())

	return engine
}

func addUniqueRequestID(c *gin.Context) {
	uniqueID := rand.New(rand.NewSource(time.Now().UnixNano()))
	c.Set("uid", strconv.Itoa(uniqueID.Int()))
}
