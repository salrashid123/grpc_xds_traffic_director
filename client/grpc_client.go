package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"

	"github.com/salrashid123/go-grpc-td/echo"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	//healthpb "google.golang.org/grpc/health/grpc_health_v1"
	_ "google.golang.org/grpc/xds" // use for xds-experimental:///be-srv
)

const ()

var ()

func main() {

	address := flag.String("host", "localhost:50051", "host:port of gRPC server")
	cacert := flag.String("cacert", "certs/tls-ca-chain.pem", "CACert for server")
	serverName := flag.String("servername", "grpc.domain.com", "CACert for server")

	flag.Parse()

	var err error

	var tlsCfg tls.Config
	rootCAs := x509.NewCertPool()
	pem, err := ioutil.ReadFile(*cacert)
	if err != nil {
		log.Fatalf("failed to load root CA certificates  error=%v", err)
	}
	if !rootCAs.AppendCertsFromPEM(pem) {
		log.Fatalf("no root CA certs parsed from file ")
	}
	tlsCfg.RootCAs = rootCAs
	tlsCfg.ServerName = *serverName

	ce := credentials.NewTLS(&tlsCfg)

	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(ce))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()

	c := echo.NewEchoServerClient(conn)

	// now make a gRPC call
	r, err := c.SayHelloUnary(ctx, &echo.EchoRequest{Name: "foo"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Unary Request Response:  %s", r.Message)
}
