package configuration

import (
	"os"
	"strconv"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// Configuration is the main structure that contain all istio-authn configuration
type Configuration struct {
	ContextName           string
	ServerAddress         string
	AccessTokenExpiration int32
	TracingEnable         bool
}

var once sync.Once

// Instance represent a single configuration instance
var Instance *Configuration

// New is a configuration factory method (sigleton)
func New() *Configuration {
	once.Do(func() {

		logrus.SetFormatter(utcFormatter{&logrus.JSONFormatter{}})
		setupConfigPath()

		Instance = new(Configuration)
		Instance.ServerAddress = viper.GetString("server.address")
		Instance.ContextName = viper.GetString("server.contextName")
		Instance.AccessTokenExpiration = int32(viper.GetInt64("server.sessionTokenExpiration"))

		if tEnabled, ok := os.LookupEnv("TRACING_ENABLED"); ok {
			tActive, _ := strconv.ParseBool(tEnabled)
			Instance.TracingEnable = tActive
		}
	})

	return Instance
}

func setupConfigPath() {
	viper.SetConfigName("config")
	configPath, exist := os.LookupEnv("CONFIG_PATH")
	if exist {
		viper.AddConfigPath(configPath)
	}
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

type utcFormatter struct {
	logrus.Formatter
}

func (u utcFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}

func printConfiguration(Instance *Configuration) {
	logrus.Infoln("Configuration Loaded")
	logrus.Infoln("=====================")
	logrus.Infoln("%+v\n", Instance)
}
