package interfaces

import (
	"github.com/kai5263499/dolon/types"
)

// Source mines packet captures and produces TCP Sessions
type Source interface {
	Pcap(string) error
	Device(string, string) error
	TcpSessionChan() chan *types.TcpSession
}
