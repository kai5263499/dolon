package dolon

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/kai5263499/dolon/generated"
)

// timeout is the length of time to wait befor flushing connections and
// bidirectional stream pairs.
const streamTimeout time.Duration = time.Minute * 5

type PacketSource struct {
	conf                   *Config
	timeout                time.Duration
	rawCompletedStreamChan chan *generated.RawCompletedStream
}

func NewPacketSource(conf *Config, rawCompletedStreamChan chan *generated.RawCompletedStream) *PacketSource {
	return &PacketSource{
		conf: conf,
		rawCompletedStreamChan: rawCompletedStreamChan,
	}
}

func (ps *PacketSource) Start() error {
	var err error

	var handle *pcap.Handle

	if len(ps.conf.PcapFile) > 0 {
		handle, err = pcap.OpenOffline(ps.conf.PcapFile)
	} else {
		handle, err = pcap.OpenLive(ps.conf.CaptureDevice, 1600, true, pcap.BlockForever)
	}

	if err != nil {
		panic(err)
	}

	logrus.Debugf("device=%s bpfFilter=%s", ps.conf.CaptureDevice, ps.conf.BPFFilter)

	if err := handle.SetBPFFilter(ps.conf.BPFFilter); err != nil {
		panic(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Set up assembly
	streamFactory := &StreamFactory{
		BidirectionalStreamMap: make(map[FlowKey]*BidirectionalStream),
		Timeout:                ps.timeout,
		RawCompletedStreamChan: ps.rawCompletedStreamChan,
	}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)
	// Limit memory usage by auto-flushing connection state if we get over 100K
	// packets in memory, or over 1000 for a single stream.
	assembler.MaxBufferedPagesTotal = 100000
	assembler.MaxBufferedPagesPerConnection = 1000

	packets := packetSource.Packets()
	ticker := time.Tick(streamTimeout / 4)
outer:
	for {
		select {
		case packet := <-packets:
			// nil packets signal end of pcap replay
			if packet == nil {
				break outer
			}
			// logrus.Infof(packet)
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				logrus.Warnf("Unusable packet")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

		case <-ticker:
			// Every minute, flush connections that haven't seen activity in the past minute.
			logrus.Debugf("---- FLUSHING ----")
			assembler.FlushOlderThan(time.Now().Add(-streamTimeout))
			streamFactory.CollectOldStreams()
		}
	}

	return nil
}
