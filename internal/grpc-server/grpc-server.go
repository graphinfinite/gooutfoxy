package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"rusprofile/config"
	"rusprofile/internal/service/rpclient"
	pb "rusprofile/pkg/grpc"
)

type RusProfileClientInterface interface {
	GetCompanyByINN(inn string) (rpclient.CompanyData, error)
}

type GrpcServer struct {
	Server           *grpc.Server
	Logger           zerolog.Logger
	RusProfileClient RusProfileClientInterface
}

func NewGrpcServer(rpc RusProfileClientInterface, log zerolog.Logger) *GrpcServer {
	server := grpc.NewServer()
	return &GrpcServer{
		Server:           server,
		RusProfileClient: rpc,
		Logger:           log,
	}
}

func httpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
		w.WriteHeader(code)
	}
	return nil
}

// запуск и инициализация grpc и gateway сервера
func (s *GrpcServer) Run(cfg config.Configuration) {
	//Grpc
	lis, err := net.Listen("tcp", net.JoinHostPort(cfg.HostGrpc, cfg.PortGrpc))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	pb.RegisterRusprofileWrapperServiceServer(s.Server, s)
	s.Logger.Info().Msgf("🚀 serving gRPC on %s", net.JoinHostPort(cfg.HostGrpc, cfg.PortGrpc))
	go func() {
		log.Fatalln(s.Server.Serve(lis))
	}()

	//Gateway
	conn, err := grpc.DialContext(
		context.Background(),
		net.JoinHostPort(cfg.HostGrpc, cfg.PortGrpc),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(httpResponseModifier),
	)
	err = pb.RegisterRusprofileWrapperServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.HostHttp, cfg.PortHttp),
		Handler: gwmux,
	}

	s.Logger.Info().Msgf("🚀 serving gateway on %s", net.JoinHostPort(cfg.HostHttp, cfg.PortHttp))
	log.Fatalln(gwServer.ListenAndServe())
}

// ping
func (s *GrpcServer) DoPing(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	s.Logger.Debug().Msg("PING")
	return &pb.PingResponse{Code: uint32(0), Message: "ok"}, nil
}

func convCompany(c rpclient.CompanyData) *pb.Company {
	return &pb.Company{Inn: c.INN, Kpp: c.KPP, Headname: c.HeadName, Name: c.Name}
}

// Получаем данные о компании по ИНН
func (s *GrpcServer) GetCompanyByINN(ctx context.Context, req *pb.GetCompanyByINNRequestV1) (*pb.GetCompanyByINNResponseV1, error) {
	s.Logger.Debug().Msgf("GetCompanyByINN: %s", req.GetInn())
	_, err := strconv.Atoi(req.Inn)
	if err != nil || len(req.GetInn()) != 10 {
		_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "400"))
		return &pb.GetCompanyByINNResponseV1{Code: uint32(0), Message: ""}, fmt.Errorf("ИНН компании должен состоять из 10 цифр")
	}
	company, err := s.RusProfileClient.GetCompanyByINN(req.GetInn())
	if err != nil {
		_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "404"))
		return &pb.GetCompanyByINNResponseV1{Code: uint32(0), Message: ""}, err
	}
	return &pb.GetCompanyByINNResponseV1{Code: uint32(0), Message: "ok", Company: convCompany(company)}, nil
}
