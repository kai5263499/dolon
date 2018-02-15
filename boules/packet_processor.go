package boules

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/kai5263499/boules/generated"
)

func NewPacketProcessor(conf *Config, rawCompletedStreamChan chan *generated.RawCompletedStream, httpStreamChan chan *generated.HttpStream) *PacketProcessor {
	return &PacketProcessor{
		conf: conf,
		rawCompletedStreamChan: rawCompletedStreamChan,
		httpStreamChan:         httpStreamChan,
	}
}

var hostRegex = regexp.MustCompile(`^[Hh]ost`)

type PacketProcessor struct {
	conf                   *Config
	rawCompletedStreamChan chan *generated.RawCompletedStream
	httpStreamChan         chan *generated.HttpStream
}

func (p *PacketProcessor) parseHeader(rawHeaderString string) (*generated.HttpHeader, error) {
	parts := strings.Split(rawHeaderString, ": ")
	if len(parts) > 1 {
		return &generated.HttpHeader{
			Key:   parts[0],
			Value: parts[1],
		}, nil
	}
	return nil, errors.New("too few header parts")
}

func (p *PacketProcessor) parseRequestPayload(rawPayloadString string) *generated.HttpRequest {
	rawPayloadString = strings.Replace(rawPayloadString, "\r", "", -1)

	payloadParts := strings.Split(rawPayloadString, "\n\n")

	if len(payloadParts) < 0 {
		return nil
	}

	headerParts := strings.Split(payloadParts[0], "\n")

	headers := make([]*generated.HttpHeader, 0)

	uriParts := strings.Split(headerParts[0], " ")

	host := ""

	for _, s := range headerParts[1:] {
		if len(s) > 1 {
			if len(host) < 1 && hostRegex.MatchString(s) {
				hostParts := strings.Split(s, ": ")
				host = hostParts[1]
				continue
			}

			header, err := p.parseHeader(s)
			if err == nil {
				headers = append(headers, header)
			}
		}
	}

	return &generated.HttpRequest{
		Host:        host,
		Uri:         uriParts[1],
		Verb:        uriParts[0],
		HttpVersion: uriParts[2],
		Headers:     headers,
		Body:        []byte(strings.Join(payloadParts[1:], "\n")),
	}
}

func (p *PacketProcessor) parseResponsePayload(rawPayloadString string) *generated.HttpResponse {
	rawPayloadString = strings.Replace(rawPayloadString, "\r", "", -1)

	payloadParts := strings.Split(rawPayloadString, "\n\n")

	if len(payloadParts) < 0 {
		return nil
	}

	headerParts := strings.Split(payloadParts[0], "\n")

	responseParts := strings.Split(headerParts[0], " ")

	code, _ := strconv.ParseInt(responseParts[1], 10, 32)

	headers := make([]*generated.HttpHeader, 0)

	for _, s := range headerParts[1:] {
		if len(s) > 1 {
			header, err := p.parseHeader(s)
			if err == nil {
				headers = append(headers, header)
			}
		}
	}

	return &generated.HttpResponse{
		Version: responseParts[0],
		Code:    int32(code),
		Headers: headers,
		Body:    []byte(strings.Join(payloadParts[1:], "\n")),
	}
}

func (p *PacketProcessor) processRawCompletedStream(rawCompletedStream *generated.RawCompletedStream) {
	spew.Dump(rawCompletedStream)

	srcString := string(rawCompletedStream.SrcData)

	httpRequest := p.parseRequestPayload(srcString)

	dstString := string(rawCompletedStream.DstData)

	httpResponse := p.parseResponsePayload(dstString)

	p.httpStreamChan <- &generated.HttpStream{
		RequestEndpoint:  rawCompletedStream.SrcEndpoint,
		Request:          httpRequest,
		ResponseEndpoint: rawCompletedStream.DstEndpoint,
		Response:         httpResponse,
	}
}

func (p *PacketProcessor) consumeRawCompletedStreamChan() {
	for rawCompletedStream := range p.rawCompletedStreamChan {
		go p.processRawCompletedStream(rawCompletedStream)
	}
}

func (p *PacketProcessor) Start() error {
	go p.consumeRawCompletedStreamChan()
	return nil
}
