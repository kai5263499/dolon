package dolon

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("source", func() {
	It("should parse a simple request from a pcap file", func() {
		var err error

		source := NewSource()

		err = source.Pcap("testdata/http.pcap")
		Expect(err).To(BeNil())

		ticker := time.NewTicker((time.Millisecond * 100))
		for {
			select {
			case evt := <-source.TcpSessionChan():
				Expect(evt.DstEndpoint.Ip).ToNot(BeNil())
				Expect(evt.DstEndpoint.Port).ToNot(BeNil())
				Expect(evt.DstData).ToNot(BeNil())
				Expect(evt.SrcEndpoint.Ip).ToNot(BeNil())
				Expect(evt.SrcEndpoint.Port).ToNot(BeNil())
				Expect(evt.SrcData).ToNot(BeNil())
			case _ = <-ticker.C:
				return
			}
		}
	})
})
