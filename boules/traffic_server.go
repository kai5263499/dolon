package boules

import (
	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/kai5263499/boules/generated"
)

func NewTrafficServer(conf *Config, httpStreamChan chan *generated.HttpStream) *TrafficServer {
	return &TrafficServer{
		conf:                  conf,
		httpStreamChan:        httpStreamChan,
		clientHttpStreamChans: make(map[generated.Traffic_GetHttpStreamServer]chan *generated.HttpStream),
	}
}

type TrafficServer struct {
	conf                  *Config
	httpStreamChan        chan *generated.HttpStream
	clientHttpStreamChans map[generated.Traffic_GetHttpStreamServer]chan *generated.HttpStream
}

func (s *TrafficServer) Start() error {

	go s.consumeHttpStreams()
	return nil
}

func (s *TrafficServer) consumeHttpStreams() {
	for httpStream := range s.httpStreamChan {
		for _, clientStreamChan := range s.clientHttpStreamChans {
			clientStreamChan <- httpStream
		}
	}
}

func (s *TrafficServer) GetHttpStream(empty *google_protobuf.Empty, stream generated.Traffic_GetHttpStreamServer) error {
	clientStreamChan := make(chan *generated.HttpStream, 0)

	s.clientHttpStreamChans[stream] = clientStreamChan

	for httpStream := range clientStreamChan {
		logrus.Infof("got httpStream")
		spew.Dump(httpStream)
		stream.Send(httpStream)
	}

	delete(s.clientHttpStreamChans, stream)

	return nil
}
