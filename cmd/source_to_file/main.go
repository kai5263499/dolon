package main

import (
	"os"
	"sync"
	"time"

	"github.com/caarlos0/env"
	"github.com/gogo/protobuf/proto"
	"github.com/kai5263499/dolon"
	"github.com/kai5263499/dolon/interfaces"
	"github.com/kai5263499/dolon/types"
	"github.com/sirupsen/logrus"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

type config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
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

func main() {
	var err error
	cfg = config{}
	err = env.Parse(&cfg)
	checkError(err)

	level, err := logrus.ParseLevel(cfg.LogLevel)
	checkError(err)

	logrus.SetLevel(level)

	parquetFilename := os.Args[1]

	fw, err := local.NewLocalFileWriter(parquetFilename)
	checkError(err)

	pw, err := writer.NewParquetWriter(fw, new(types.ProtoRecord), 4)
	checkError(err)

	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	source := dolon.NewSource()

	var wg sync.WaitGroup

	for _, arg := range os.Args[2:] {
		wg.Add(1)
		go processPcap(arg, source, pw, &wg)
	}

	wg.Wait()

	err = pw.WriteStop()
	checkError(err)

	fw.Close()
}

func processPcap(filename string, source interfaces.Source, pw *writer.ParquetWriter, wg *sync.WaitGroup) {
	defer wg.Done()

	var err error

	err = source.Pcap(filename, "")
	checkError(err)

	ticker := time.NewTicker((time.Millisecond * 100))

	for {
		select {
		case evt := <-source.TcpSessionChan():
			data, err := proto.Marshal(evt)
			checkError(err)

			rec := types.ProtoRecord{
				Data: string(data),
			}

			err = pw.Write(rec)
			checkError(err)
		case _ = <-ticker.C:
			return
		}
	}
}
