package dolon

import (
	"fmt"

	"github.com/google/gopacket"
)

// FlowKey is used to map bidirectional streams to each other
type FlowKey struct {
	net, transport gopacket.Flow
}

// String prints out the key in a human-readable fashion.
func (k *FlowKey) String() string {
	return fmt.Sprintf("%v:%v", k.net, k.transport)
}
