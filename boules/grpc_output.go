package boules

import (
	"fmt"
	"log"
	"net"

	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/kai5263499/boules/generated"
)

func NewGrpcOutput(conf *Config, trafficServer *TrafficServer) *GrpcOutput {
	return &GrpcOutput{
		conf:          conf,
		trafficServer: trafficServer,
	}
}

type GrpcOutput struct {
	conf          *Config
	trafficServer *TrafficServer
	lis           *net.Listener
	grpcServer    *grpc.Server
}

func (s *GrpcOutput) Start() error {
	var err error

	listenAddress := fmt.Sprintf("localhost:%d", s.conf.GrpcPort)

	lis, err := net.Listen("tcp", listenAddress)

	if err != nil {
		logrus.Errorf("unable to create network listener err=%#v", err)
		return err
	}

	logrus.Infof("listening on %s", listenAddress)

	s.lis = &lis

	if !s.conf.UseTLS {
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			return err
		}
		s.grpcServer = grpc.NewServer()
	} else {
		certFile := s.conf.SSLCertFile
		keyFile := s.conf.SSLKeyFile
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			return err
		}
		s.grpcServer = grpc.NewServer(grpc.Creds(creds))
	}

	generated.RegisterTrafficServer(s.grpcServer, s.trafficServer)

	return s.grpcServer.Serve(lis)
}
