package boules

import (
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/kai5263499/boules/generated"
)

func NewTrafficServer(conf *Config, httpStreamChan chan *generated.HttpStream) *TrafficServer {
	return &TrafficServer{
		conf:           conf,
		httpStreamChan: httpStreamChan,
	}
}

type TrafficServer struct {
	conf           *Config
	httpStreamChan chan *generated.HttpStream
}

func (s *TrafficServer) Start() error {
	return nil
}

func (s *TrafficServer) GetHttpStream(empty *google_protobuf.Empty, stream generated.Traffic_GetHttpStreamServer) error {

	return nil
}
