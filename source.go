package dolon

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/kai5263499/dolon/interfaces"
	"github.com/kai5263499/dolon/internal"
	"github.com/kai5263499/dolon/types"
	"github.com/sirupsen/logrus"
)

var _ interfaces.Source = (*Source)(nil)

func NewSource() *Source {
	return &Source{
		tcpSessionChan: make(chan *types.TcpSession, 1000),
	}
}

// timeout is the length of time to wait befor flushing connections and
// bidirectional stream pairs.
const streamTimeout time.Duration = time.Minute * 2

type Source struct {
	tcpSessionChan chan *types.TcpSession
	timeout        time.Duration
}

func (s *Source) Pcap(pcapFile string) error {
	var err error

	logrus.WithFields(logrus.Fields{
		"pcapFile": pcapFile,
	}).Infof("starting process packet loop")

	handle, err := pcap.OpenOffline(pcapFile)
	if err != nil {
		return err
	}

	s.timeout = time.Millisecond * 100

	go s.processPackets(handle)

	return nil
}

func (s *Source) Device(captureDevice, bpfFilter string) error {
	var err error

	logrus.WithFields(logrus.Fields{
		"captureDevice": captureDevice,
		"bpfFilter":     bpfFilter,
	}).Infof("starting process packet loop")

	handle, err := pcap.OpenLive(captureDevice, 65536, true, pcap.BlockForever)
	if err != nil {
		return err
	}

	err = handle.SetBPFFilter(bpfFilter)
	if err != nil {
		return err
	}

	s.timeout = time.Minute * 5

	go s.processPackets(handle)

	return nil
}

func (s *Source) processPackets(handle *pcap.Handle) {
	source := gopacket.NewPacketSource(handle, handle.LinkType())

	// Set up assembly
	streamFactory := &internal.StreamFactory{
		BidirectionalStreamMap: make(map[internal.FlowKey]*internal.BidirectionalStream),
		Timeout:                s.timeout,
		TcpSessionChan:         s.tcpSessionChan,
	}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)
	// Limit memory usage by auto-flushing connection state if we get over 100K
	// packets in memory, or over 1000 for a single stream.
	assembler.MaxBufferedPagesTotal = 100000
	assembler.MaxBufferedPagesPerConnection = 1000

	packets := source.Packets()
	ticker := time.NewTicker(s.timeout)
	for {
		select {
		case packet := <-packets:
			// nil packets signal end of pcap replay
			if packet == nil {
				assembler.FlushOlderThan(time.Now().Add(-streamTimeout))
				streamFactory.CollectOldStreams()

				return
			}
			// logrus.Infof(packet)
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		case _ = <-ticker.C:
			// Every minute, flush connections that haven't seen activity in the past minute.
			assembler.FlushOlderThan(time.Now().Add(-streamTimeout))
			streamFactory.CollectOldStreams()
		}
	}
}

func (s *Source) TcpSessionChan() chan *types.TcpSession {
	return s.tcpSessionChan
}
