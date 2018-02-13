package boules

import (
	"encoding/json"
	"fmt"

	"github.com/kai5263499/boules/generated"
)

func NewConsoleOutput(conf *Config, httpStreamChan chan *generated.HttpStream) *ConsoleOutput {
	return &ConsoleOutput{
		conf:           conf,
		httpStreamChan: httpStreamChan,
	}
}

type ConsoleOutput struct {
	conf           *Config
	httpStreamChan chan *generated.HttpStream
}

func (o *ConsoleOutput) consumeHttpStreamChan() {
	for httpStream := range o.httpStreamChan {
		switch o.conf.OutputFormat {
		case CurlifyOutputFormat:
			fmt.Println(Curlify(httpStream))
		case JsonOutputFormat:
			jsonString, _ := json.Marshal(httpStream)
			fmt.Println(string(jsonString))
		}
	}
}

func (o *ConsoleOutput) Start() error {
	go o.consumeHttpStreamChan()
	return nil
}
