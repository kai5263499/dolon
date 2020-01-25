package internal

import (
	"strconv"
	"time"

	"github.com/kai5263499/dolon/types"
	"github.com/sirupsen/logrus"
)

// bidi stores each unidirectional side of a bidirectional stream.
//
// When a new stream comes in, if we don't have an opposite stream, a bidi is
// created with 'a' set to the new stream.  If we DO have an opposite stream,
// 'b' is set to the new stream.
type BidirectionalStream struct {
	key            FlowKey                // Key of the first stream, mostly for logging.
	a, b           *Stream                // the two bidirectional streams.
	lastPacketSeen time.Time              // last time we saw a packet from either stream.
	tcpSessionChan chan *types.TcpSession // Channel to output raw completed streams to
}

// maybeFinish will wait until both directions are complete, then print out
// stats.
func (bd *BidirectionalStream) maybeFinish() {
	switch {
	case bd.a == nil:
		logrus.Errorf("[%v] a should always be non-nil, since it's set when bidis are created", bd.key)
	case !bd.a.done:
		logrus.Tracef("[%v] still waiting on first stream", bd.key)
	case bd.b == nil:
		logrus.Tracef("[%v] no second stream yet", bd.key)
	case !bd.b.done:
		logrus.Tracef("[%v] still waiting on second stream", bd.key)
	default:
		if bd.a.bytes > 0 && bd.b.bytes > 0 {

			srcPort, err := strconv.ParseInt(bd.key.transport.Src().String(), 10, 32)
			if err != nil {
				return
			}

			dstPort, err := strconv.ParseInt(bd.key.transport.Dst().String(), 10, 32)
			if err != nil {
				return
			}

			rs := &types.TcpSession{
				SrcEndpoint: &types.Endpoint{
					Ip:   bd.key.net.Src().String(),
					Port: int32(srcPort),
				},
				SrcData: bd.a.completePayload.Bytes(),
				DstEndpoint: &types.Endpoint{
					Ip:   bd.key.net.Dst().String(),
					Port: int32(dstPort),
				},
				DstData: bd.b.completePayload.Bytes(),
			}
			bd.tcpSessionChan <- rs
		}
	}
}
