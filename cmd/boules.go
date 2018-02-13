package main

import (
	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/kai5263499/boules/boules"
	"github.com/kai5263499/boules/generated"
	"github.com/oklog/run"
)

var (
	device       string
	bpfFilter    string
	pcapFile     string
	outputType   string
	outputFormat string
	grpcPort     int
	useTLS       bool
	sslCertFile  string
	sslKeyFile   string
)

func main() {
	flag.StringVar(&pcapFile, "pcap", "", "pcap file to read from (disables live mode)")
	flag.StringVar(&device, "device", "eth0", "device")
	flag.StringVar(&bpfFilter, "filter", "tcp and dst port 8080", "filter")
	flag.StringVar(&outputType, "outputType", "console", "how output should be handled (grpc, console)")
	flag.StringVar(&outputFormat, "outputFormat", "curlify", "how output should be formatted (curlify, json)")
	flag.IntVar(&grpcPort, "grpcPort", 9001, "")
	flag.BoolVar(&useTLS, "useTLS", false, "")
	flag.StringVar(&sslCertFile, "sslCert", "ssl.crt", "")
	flag.StringVar(&sslKeyFile, "sslKey", "ssl.key", "")

	flag.Parse()

	logrus.SetLevel(logrus.InfoLevel)

	conf := &boules.Config{
		CaptureDevice: device,
		BPFFilter:     bpfFilter,
		PcapFile:      pcapFile,
		OutputType:    boules.OutputType(outputType),
		OutputFormat:  boules.OutputFormat(outputFormat),
		GrpcPort:      grpcPort,
		UseTLS:        useTLS,
		SSLCertFile:   sslCertFile,
		SSLKeyFile:    sslKeyFile,
	}

	rawCompletedStreamChan := make(chan *generated.RawCompletedStream, 1000)
	httpCompletedStreamChan := make(chan *generated.HttpStream, 1000)

	packetSource := boules.NewPacketSource(conf, rawCompletedStreamChan)

	packetProcessor := boules.NewPacketProcessor(conf, rawCompletedStreamChan, httpCompletedStreamChan)

	var g run.Group
	g.Add(func() error {
		return packetSource.Start()
	}, func(error) {
	})

	g.Add(func() error {
		return packetProcessor.Start()
	}, func(error) {
	})

	switch conf.OutputType {
	case boules.ConsoleOutputType:
		consoleOutput := boules.NewConsoleOutput(conf, httpCompletedStreamChan)
		g.Add(func() error {
			return consoleOutput.Start()
		}, func(error) {
		})
	case boules.GrpcOutputType:
		grpcOutput := boules.NewGrpcOutput(conf, httpCompletedStreamChan)
		g.Add(func() error {
			return grpcOutput.Start()
		}, func(error) {
		})
	}
	g.Run()
}
