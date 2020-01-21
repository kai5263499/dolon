package main

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"

	"github.com/kai5263499/dolon"
	"github.com/kai5263499/dolon/interfaces"
	"github.com/kai5263499/dolon/types"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

type config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

var (
	cfg config
)

func checkError(err error) {
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Errorf("encountered unrecoverable err")
		panic(err)
	}
}

var _ interfaces.Source = (*mockSource)(nil)

var (
	tcpSessionChan chan *types.TcpSession
)

type mockSource struct{}

func (s *mockSource) Pcap(pcapFile string) error {
	return nil
}

func (s *mockSource) Device(captureDevice, bpfFilter string) error {
	return nil
}

func (s *mockSource) TcpSessionChan() chan *types.TcpSession {
	return tcpSessionChan
}

func main() {
	var err error
	cfg = config{}
	err = env.Parse(&cfg)
	checkError(err)

	parquetFilename := os.Args[1]

	tcpSessionChan = make(chan *types.TcpSession, 100)

	ms := mockSource{}

	processor := dolon.NewProcessor()
	processor.StartProcessingTcpSessions(&ms)

	fr, err := local.NewLocalFileReader(parquetFilename)
	checkError(err)

	pr, err := reader.NewParquetReader(fr, new(types.ProtoRecord), 4)
	checkError(err)

	num := int(pr.GetNumRows())
	checkError(err)

	protoRecords := make([]types.ProtoRecord, num)
	err = pr.Read(&protoRecords)
	checkError(err)

	for _, protoRecord := range protoRecords {
		var tcpSession types.TcpSession
		err = proto.Unmarshal([]byte(protoRecord.Data), &tcpSession)
		checkError(err)

		tcpSessionChan <- &tcpSession
	}

	pr.ReadStop()
	fr.Close()

	ticker := time.NewTicker((time.Millisecond * 100))
	for {
		select {
		case httpSession := <-processor.HttpSessionChan():
			fmt.Printf("%s\n", dolon.Curlify(httpSession))
		case _ = <-ticker.C:
			return
		}
	}
}
