package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
	"github.com/Sirupsen/logrus"
	"fmt"
	"flag"
)

var (
	port uint
)

func main() {

	// TODO: Make device configurable

	flag.UintVar(&port, "port", 80, "port")

	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)

	handle, err := pcap.OpenLive("en0", 1600, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}

	captureCmd := fmt.Sprintf("tcp and dst port %d", port)
	captureCmd = "arp"

	logrus.Debugf("captureCmd=%s", captureCmd)

	if err := handle.SetBPFFilter(captureCmd); err != nil {
		panic(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		logrus.Debugf("Packet data: packet=%#v", packet)
		for _, layer := range packet.Layers() {
			if layer.LayerType() == layers.LayerTypeTCP && len(layer.LayerPayload()) > 0{
				logrus.Infof("TCP layer payload: %#v", string(layer.LayerPayload()))
			}
		}
	}
}