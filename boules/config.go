package boules

type Config struct {
	CaptureDevice string
	BPFFilter     string
	PcapFile      string
	OutputType    string
	GrpcPort      int
	UseTLS        bool
	SSLCertFile   string
	SSLKeyFile    string
}
