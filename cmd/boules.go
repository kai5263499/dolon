package main

import (
	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device    string
	bpfFilter string
)

func main() {

	flag.StringVar(&device, "device", "eth0", "device")
	flag.StringVar(&bpfFilter, "filter", "tcp and dst port ", "filter")

	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)

	handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}

	logrus.Debugf("device=%s bpfFilter=%s", device, bpfFilter)

	if err := handle.SetBPFFilter(bpfFilter); err != nil {
		panic(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		for _, layer := range packet.Layers() {
			if layer.LayerType() == layers.LayerTypeTCP {
				layerPayload := string(layer.LayerPayload())
				logrus.Infof("TCP layer payload: [%s]", layerPayload)
			}
		}
	}
}
