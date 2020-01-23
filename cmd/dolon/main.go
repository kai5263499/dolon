package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/kai5263499/dolon"
	"github.com/kai5263499/dolon/interfaces"
	"github.com/sirupsen/logrus"
)

type config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
}

var (
	cfg       config
	device    string
	bpfFilter string
)

func checkError(err error) {
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatalf("encountered unrecoverable err")
	}
}

func main() {
	var err error
	cfg = config{}
	err = env.Parse(&cfg)
	checkError(err)

	flag.StringVar(&device, "device", "lo", "device")
	flag.StringVar(&bpfFilter, "filter", "tcp and dst port 8080", "filter")

	flag.Parse()

	level, err := logrus.ParseLevel(cfg.LogLevel)
	checkError(err)

	logrus.SetLevel(level)

	source := dolon.NewSource()
	var wg sync.WaitGroup

	wg.Add(1)
	go processFromDevice(source, &wg)

	processor := dolon.NewProcessor()
	processor.StartProcessingTcpSessions(source)

	wg.Add(1)
	go processHttpSessions(processor, &wg)

	wg.Wait()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func processFromDevice(source interfaces.Source, wg *sync.WaitGroup) {
	defer wg.Done()

	logrus.Infof("before device")
	err := source.Device(device, bpfFilter)
	checkError(err)
}

func processHttpSessions(processor interfaces.HttpProcessor, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker((time.Minute * 5))
	for {
		select {
		case httpSession := <-processor.HttpSessionChan():
			fmt.Printf("%s\n", dolon.Curlify(httpSession))
		case _ = <-ticker.C:
			return
		}
	}
}
