package dolon

import (
	"bytes"

	"github.com/google/gopacket/tcpassembly"
)

// Stream implements tcpassembly.Stream
type Stream struct {
	bytes               int64 // total bytes seen on this stream.
	completePayload     *bytes.Buffer
	bidirectionalStream *BidirectionalStream // maps to my bidirectional twin.
	done                bool                 // if true, we've seen the last packet we're going to for this stream.
}

// EmptyStream is used to finish bidi that only have one stream, in
// collectOldStreams.
var EmptyStream = &Stream{done: true}

// Reassembled handles reassembled TCP stream data.
func (s *Stream) Reassembled(rs []tcpassembly.Reassembly) {
	for _, r := range rs {
		// For now, we'll simply count the bytes on each side of the TCP stream.
		s.bytes += int64(len(r.Bytes))
		if r.Skip > 0 {
			s.bytes += int64(r.Skip)
		}
		s.completePayload.Write(r.Bytes)

		// Mark that we've received new packet data.
		// We could just use time.Now, but by using r.Seen we handle the case
		// where packets are being read from a file and could be very old.
		if s.bidirectionalStream.lastPacketSeen.After(r.Seen) {
			s.bidirectionalStream.lastPacketSeen = r.Seen
		}
	}
}

// ReassemblyComplete marks this stream as finished.
func (s *Stream) ReassemblyComplete() {
	s.done = true
	s.bidirectionalStream.maybeFinish()
}
