package main

import (
	"flag"
	"log"
	"net"

	"github.com/salrashid123/go-grpc-td/echo"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	tlsCert   = flag.String("tlsCert", "certs/grpc.crt", "tls Certificate")
	tlsKey    = flag.String("tlsKey", "certs/grpc.key", "tls Key")
	grpcport  = flag.String("grpcport", ":50051", "grpcport")
	hcaddress = flag.String("hcaddress", ":50050", "host:port of gRPC HC server")
)

const ()

type server struct {
	echo.UnimplementedEchoServerServer
}

type healthServer struct {
}

func (s *server) SayHelloUnary(ctx context.Context, in *echo.EchoRequest) (*echo.EchoReply, error) {
	log.Println("Got Unary Request: ")
	uid, _ := uuid.NewUUID()
	return &echo.EchoReply{Message: "SayHelloUnary Response " + uid.String()}, nil
}

func (s *healthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	log.Printf("HealthCheck Called for Service %s", in.Service)
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

func main() {

	flag.Parse()

	go func() {
		lis, err := net.Listen("tcp", *hcaddress)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}

		s := grpc.NewServer(sopts...)
		healthpb.RegisterHealthServer(s, &healthServer{})
		s.Serve(lis)

	}()

	lis, err := net.Listen("tcp", *grpcport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}

	ce, err := credentials.NewServerTLSFromFile(*tlsCert, *tlsKey)
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}
	sopts = append(sopts, grpc.Creds(ce))

	s := grpc.NewServer(sopts...)
	echo.RegisterEchoServerServer(s, &server{})

	log.Println("Starting server")
	s.Serve(lis)

}
