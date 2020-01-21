package dolon

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/kai5263499/dolon/interfaces"
	"github.com/kai5263499/dolon/types"
)

func NewProcessor() *PacketProcessor {
	return &PacketProcessor{
		httpStreamChan: make(chan *types.HttpSession, 0),
	}
}

var _ interfaces.HttpProcessor = (*PacketProcessor)(nil)

var hostRegex = regexp.MustCompile(`^[Hh]ost`)

type PacketProcessor struct {
	tcpSessionChan chan *types.TcpSession
	httpStreamChan chan *types.HttpSession
}

func (p *PacketProcessor) parseHeader(rawHeaderString string) (*types.HttpHeader, error) {
	parts := strings.Split(rawHeaderString, ": ")
	if len(parts) > 1 {
		return &types.HttpHeader{
			Key:   parts[0],
			Value: parts[1],
		}, nil
	}
	return nil, errors.New("too few header parts")
}

func (p *PacketProcessor) parseRequestPayload(rawPayloadString string) *types.HttpRequest {
	rawPayloadString = strings.Replace(rawPayloadString, "\r", "", -1)

	payloadParts := strings.Split(rawPayloadString, "\n\n")

	if len(payloadParts) < 0 {
		return nil
	}

	headerParts := strings.Split(payloadParts[0], "\n")

	headers := make([]*types.HttpHeader, 0)

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
			if err != nil {
				continue
			}
			headers = append(headers, header)
		}
	}

	return &types.HttpRequest{
		Host:        host,
		Uri:         uriParts[1],
		Verb:        uriParts[0],
		HttpVersion: uriParts[2],
		Headers:     headers,
		Body:        []byte(strings.Join(payloadParts[1:], "\n")),
	}
}

func (p *PacketProcessor) parseResponsePayload(rawPayloadString string) *types.HttpResponse {
	rawPayloadString = strings.Replace(rawPayloadString, "\r", "", -1)

	payloadParts := strings.Split(rawPayloadString, "\n\n")

	if len(payloadParts) < 0 {
		return nil
	}

	headerParts := strings.Split(payloadParts[0], "\n")

	responseParts := strings.Split(headerParts[0], " ")

	code, _ := strconv.ParseInt(responseParts[1], 10, 32)

	headers := make([]*types.HttpHeader, 0)

	for _, s := range headerParts[1:] {
		if len(s) > 1 {
			header, err := p.parseHeader(s)
			if err == nil {
				headers = append(headers, header)
			}
		}
	}

	return &types.HttpResponse{
		Version: responseParts[0],
		Code:    int32(code),
		Headers: headers,
		Body:    []byte(strings.Join(payloadParts[1:], "\n")),
	}
}

func (p *PacketProcessor) processRawCompletedStream(rawCompletedStream *types.TcpSession) {
	srcString := string(rawCompletedStream.SrcData)

	httpRequest := p.parseRequestPayload(srcString)

	dstString := string(rawCompletedStream.DstData)

	httpResponse := p.parseResponsePayload(dstString)

	p.httpStreamChan <- &types.HttpSession{
		Request:  httpRequest,
		Response: httpResponse,
	}
}

func (p *PacketProcessor) StartProcessingTcpSessions(source interfaces.Source) {
	go p.processLoop(source)
}

func (p *PacketProcessor) processLoop(source interfaces.Source) {
	for {
		select {
		case evt := <-source.TcpSessionChan():
			p.processRawCompletedStream(evt)
		}
	}
}

func (p *PacketProcessor) HttpSessionChan() chan *types.HttpSession {
	return p.httpStreamChan
}
