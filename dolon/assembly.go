package dolon

import (
	"bytes"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/kai5263499/dolon/generated"
)

// StreamFactory implements tcpassmebly.StreamFactory
type StreamFactory struct {
	// bidirectionalStreamMap maps keys to bidirectional stream pairs.
	BidirectionalStreamMap map[FlowKey]*BidirectionalStream
	Timeout                time.Duration
	RawCompletedStreamChan chan *generated.RawCompletedStream
}

// New handles creating a new tcpassembly.Stream.
func (f *StreamFactory) New(netFlow, tcpFlow gopacket.Flow) tcpassembly.Stream {
	// Create a new stream.
	s := &Stream{completePayload: &bytes.Buffer{}}

	// Find the bidi bidirectional struct for this stream, creating a new one if
	// one doesn't already exist in the map.
	k := FlowKey{netFlow, tcpFlow}
	bd := f.BidirectionalStreamMap[k]
	if bd == nil {
		bd = &BidirectionalStream{a: s, key: k, rawCompletedStreamChan: f.RawCompletedStreamChan}
		logrus.Debugf("[%v] created first side of bidirectional stream", bd.key)
		// Register bidirectional with the reverse key, so the matching stream going
		// the other direction will find it.
		f.BidirectionalStreamMap[FlowKey{netFlow.Reverse(), tcpFlow.Reverse()}] = bd
	} else {
		logrus.Debugf("[%v] found second side of bidirectional stream", bd.key)
		bd.b = s
		// Clear out the bidi we're using from the map, just in case.
		delete(f.BidirectionalStreamMap, k)
	}
	s.bidirectionalStream = bd
	return s
}

// CollectOldStreams finds any streams that haven't received a packet within
// 'timeout', and sets/finishes the 'b' stream inside them.  The 'a' stream may
// still receive packets after this.
func (f *StreamFactory) CollectOldStreams() {
	cutoff := time.Now().Add(-f.Timeout)
	for k, bd := range f.BidirectionalStreamMap {
		if bd.lastPacketSeen.Before(cutoff) {
			logrus.Warnf("[%v] timing out old stream", bd.key)
			bd.b = EmptyStream                  // stub out b with an empty stream.
			delete(f.BidirectionalStreamMap, k) // remove it from our map.
			bd.maybeFinish()                    // if b was the last stream we were waiting for, finish up.
		}
	}
}
