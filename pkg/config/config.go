/*
 * @File: config.config.go
 * @Description: Defines common service configuration
 * @Author: Yoan Yomba (yoanyombapro@gmail.com)
 */
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

)

// Configuration stores setting values
type Configuration struct {
	Debug       string `json:"debug.addr"`
	Http        string `json:"http.addr"`
	Appdash     string `json:"appdash.addr"`
	ZipkinUrl   string `json:"zipkin.url"`
	UseZipkin   bool   `json:"zipkin.use"`
	Zipkin      string `json:"zipkin.addr"`
	DbType      string `json:"dbType"`
	DbAddress   string `json:"dbAddress"`
	DbName      string `json:"dbName"`
	DbSettings  string `json:"dbSettings"`
	Development bool   `json:"development"`
	Jwt         string `json:"jwtSecretPassword"`
	Issuer      string `json:"issuer"`
	ServiceName string `json:"serviceName"`
}

// Config shares the global configuration
var (
	Config *Configuration
)

// Status Text
const (
	ErrNameEmpty      = "Name is empty"
	ErrPasswordEmpty  = "Password is empty"
	ErrNotObjectIDHex = "String is not a valid hex representation of an ObjectId"
)

// Status Code
const (
	StatusCodeUnknown = -1
	StatusCodeOK      = 1000
)

// LoadConfig loads configuration from the config file
func LoadConfig() {
	// Filename is the path to the json config file
	Config = new(Configuration)
	Config.ServiceName = "users_microservice"
	Config.Issuer = "lensplatform"
	Config.Jwt = "lensplatformjwtpassword"
	Config.Development = false
	Config.DbSettings = "?sslmode=require"
	Config.DbName = "users-microservice-db"
	Config.DbAddress = "doadmin:x9nec6ffkm1i3187@backend-datastore-do-user-6612421-0.db.ondigitalocean.com:25060/"
	Config.DbType = "postgresql://"
	Config.Zipkin = ":8080"
	Config.UseZipkin = true
	Config.ZipkinUrl = "http://localhost:9411/api/v2/spans"
	Config.Appdash = ":8086"
	Config.Http = ":8085"
	Config.Debug = ":8084"
}

func Init(){
	err := godotenv.Load("development.env")
	if err != nil {
		fmt.Println("Unable to parse env file")
	}

	debug := os.Getenv("debug.addr")
	fmt.Println("debug", debug)

}

// GetDatabaseConnectionString Creates a database connection string from the service configuration settings
func (Config *Configuration) GetDatabaseConnectionString() string {
	if Config == nil {
		LoadConfig()
	}
	return Config.DbType + Config.DbAddress + Config.DbName + Config.DbSettings
}
