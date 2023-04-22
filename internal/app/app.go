package app

import (
	"os"
	"rusprofile/config"
	grpcserver "rusprofile/internal/grpc-server"
	"rusprofile/internal/service/rpclient"
	"time"

	"github.com/rs/zerolog"
)

func Run() {
	cfg := config.GetConfig()
	zlog := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	rpc := rpclient.NewClient("https://www.rusprofile.ru", time.Minute)

	//TODO add usecase layeer

	server := grpcserver.NewGrpcServer(rpc, &zlog)
	server.Run(cfg)

}
