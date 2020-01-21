package dolon

import (
	"github.com/gogo/protobuf/proto"
	gomock "github.com/golang/mock/gomock"
	"github.com/kai5263499/dolon/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

var _ = Describe("processor", func() {
	var (
		mockCtrl *gomock.Controller
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should parse an http session from a stream", func() {
		tcpSessionChan := make(chan *types.TcpSession, 10)

		mockSource := NewMockSource(mockCtrl)
		mockSource.EXPECT().TcpSessionChan().AnyTimes().Return(tcpSessionChan)

		fr, err := local.NewLocalFileReader("testdata/http.parquet")
		Expect(err).To(BeNil())

		pr, err := reader.NewParquetReader(fr, new(types.ProtoRecord), 4)
		Expect(err).To(BeNil())

		num := int(pr.GetNumRows())
		Expect(num).To(Equal(2))

		protoRecords := make([]types.ProtoRecord, 2)
		err = pr.Read(&protoRecords)
		Expect(err).To(BeNil())

		for _, protoRecord := range protoRecords {
			var tcpSession types.TcpSession
			err = proto.Unmarshal([]byte(protoRecord.Data), &tcpSession)
			Expect(err).To(BeNil())

			tcpSessionChan <- &tcpSession
		}

		pr.ReadStop()
		fr.Close()

		processor := NewProcessor()
		processor.StartProcessingTcpSessions(mockSource)
		httpStreamChan := processor.HttpSessionChan()

		httpStream := <-httpStreamChan
		Expect(httpStream.Request.GetHost()).To(Equal("www.ethereal.com"))
		Expect(httpStream.Request.GetVerb()).To(Equal("GET"))

		httpStream = <-httpStreamChan
		Expect(httpStream.Request.GetHost()).To(Equal("pagead2.googlesyndication.com"))
		Expect(httpStream.Request.GetVerb()).To(Equal("GET"))
	})
})
