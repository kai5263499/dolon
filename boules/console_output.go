package boules

import "github.com/kai5263499/boules/generated"

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

func (o ConsonleOutput) consumeHttpStreamChan() {
	for httpStream := range o.consumeHttpStreamChan {

	}
}

func (o *ConsoleOutput) Start() error {
	go consumeHttpStreamChan()
	return nil
}
