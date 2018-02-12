package boules

import (
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kai5263499/boules/generated"
)

// bidi stores each unidirectional side of a bidirectional stream.
//
// When a new stream comes in, if we don't have an opposite stream, a bidi is
// created with 'a' set to the new stream.  If we DO have an opposite stream,
// 'b' is set to the new stream.
type BidirectionalStream struct {
	key                    FlowKey                            // Key of the first stream, mostly for logging.
	a, b                   *Stream                            // the two bidirectional streams.
	lastPacketSeen         time.Time                          // last time we saw a packet from either stream.
	rawCompletedStreamChan chan *generated.RawCompletedStream // Channel to output raw completed streams to
}

// maybeFinish will wait until both directions are complete, then print out
// stats.
func (bd *BidirectionalStream) maybeFinish() {
	switch {
	case bd.a == nil:
		logrus.Errorf("[%v] a should always be non-nil, since it's set when bidis are created", bd.key)
	case !bd.a.done:
		logrus.Debugf("[%v] still waiting on first stream", bd.key)
	case bd.b == nil:
		logrus.Debugf("[%v] no second stream yet", bd.key)
	case !bd.b.done:
		logrus.Debugf("[%v] still waiting on second stream", bd.key)
	default:
		if bd.a.bytes > 0 && bd.b.bytes > 0 {
			if bd.b.bytes > 1000 || bd.b.bytes == 0 {
				return
			}

			srcPort, _ := strconv.ParseInt(bd.key.transport.Src().String(), 10, 32)
			dstPort, _ := strconv.ParseInt(bd.key.transport.Dst().String(), 10, 32)

			rs := &generated.RawCompletedStream{
				SrcEndpoint: &generated.Endpoint{
					Ip:   bd.key.net.Src().String(),
					Port: int32(srcPort),
				},
				SrcData: bd.a.completePayload.Bytes(),
				DstEndpoint: &generated.Endpoint{
					Ip:   bd.key.net.Dst().String(),
					Port: int32(dstPort),
				},
				DstData: bd.b.completePayload.Bytes(),
			}
			bd.rawCompletedStreamChan <- rs

			// logrus.Infof("[%v] FINISHED, bytes: %d tx, %d rx\n", bd.key, bd.a.bytes, bd.b.bytes)
			// logrus.Infof("tx payload: %#v\n", string(bd.a.completePayload.Bytes()))
			// if bd.b.bytes < 1000 {
			// 	logrus.Infof("rx payload: %#v\n", string(bd.b.completePayload.Bytes()))
			// }
		}
	}
}
