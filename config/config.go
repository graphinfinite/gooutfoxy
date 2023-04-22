package config

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/spf13/viper"
)

/*
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
*/

// чтение конфигуарционного файла и переменных среды окружения
type Configuration struct {
	PortHttp string `json:"portHttp"`
	PortGrpc string `json:"portGrpc"`
	HostHttp string `json:"hostHttp"`
	HostGrpc string `json:"hostGrpc"`
}

var Config Configuration

// инициализация
func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("server.hostGrpc", "0.0.0.0")
	viper.SetDefault("server.hostHttp", "0.0.0.0")
	viper.SetDefault("server.portGrpc", "80")
	viper.SetDefault("server.portHttp", "81")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Config file not found \n", err)
	}
	Config = Configuration{
		PortGrpc: viper.GetString("server.portGrpc"),
		PortHttp: viper.GetString("server.portHttp"),
		HostHttp: viper.GetString("server.hostHttp"),
		HostGrpc: viper.GetString("server.hostGrpc"),
	}
}

func GetConfig() Configuration {
	return Config
}

func (config Configuration) String() string {
	out, err := json.Marshal(config)
	if err != nil {
		return ""
	}
	return (string(out))
}
