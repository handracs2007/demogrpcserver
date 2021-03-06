package main

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "github.com/handracs2007/demogrpcserver/rpc"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "io/ioutil"
    "log"
    "net"
)

type GrpcDemoService struct {
    rpc.UnimplementedDemoServiceServer
}

func (service GrpcDemoService) SayHello(context context.Context, request *rpc.HelloRequest) (*rpc.HelloResponse, error) {
    name := request.Name
    age := request.Age
    resp := fmt.Sprintf("Hello %s, you are %d year(s) old.", name, age)

    response := &rpc.HelloResponse{Response: resp}
    return response, nil
}

func main() {
    // Load the server certificate and its key
    serverCert, err := tls.LoadX509KeyPair("server.pem", "server.key")
    if err != nil {
        log.Fatalf("Failed to load server certificate and key. %s.", err)
    }

    // Load the CA certificate
    trustedCert, err := ioutil.ReadFile("cacert.pem")
    if err != nil {
        log.Fatalf("Failed to load trusted certificate. %s.", err)
    }

    // Put the CA certificate to certificate pool
    certPool := x509.NewCertPool()
    if !certPool.AppendCertsFromPEM(trustedCert) {
        log.Fatalf("Failed to append trusted certificate to certificate pool. %s.", err)
    }

    // Create the TLS configuration
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{serverCert},
        RootCAs:      certPool,
        ClientCAs:    certPool,
        MinVersion:   tls.VersionTLS13,
        MaxVersion:   tls.VersionTLS13,
    }

    // Create a new TLS credentials based on the TLS configuration
    cred := credentials.NewTLS(tlsConfig)

    // Create a listener that listens to localhost port 8443
    listener, err := net.Listen("tcp", "localhost:8443")
    if err != nil {
        log.Fatalf("Failed to start listener. %s.", err)
    }
    defer func() {
        err = listener.Close()
        if err != nil {
            log.Printf("Failed to close listener. %s\n", err)
        }
    }()

    // Create a new gRPC server
    server := grpc.NewServer(grpc.Creds(cred))
    rpc.RegisterDemoServiceServer(server, &GrpcDemoService{}) // Register the demo service

    // Start the gRPC server
    err = server.Serve(listener)
    if err != nil {
        log.Fatalf("Failed to start gRPC server. %s.", err)
    }
}
