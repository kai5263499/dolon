package boules

type OutputType string

const (
	GrpcOutputType    OutputType = "grpc"
	ConsoleOutputType OutputType = "console"
)

type OutputFormat string

const (
	CurlifyOutputFormat OutputFormat = "curlify"
	JsonOutputFormat    OutputFormat = "json"
)

type Config struct {
	CaptureDevice string
	BPFFilter     string
	PcapFile      string
	OutputType    OutputType
	OutputFormat  OutputFormat
	GrpcPort      int
	UseTLS        bool
	SSLCertFile   string
	SSLKeyFile    string
}
