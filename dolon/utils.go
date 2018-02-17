package dolon

import (
	"bytes"
	"fmt"

	"github.com/kai5263499/dolon/generated"
)

func Curlify(httpStream *generated.HttpStream) string {
	var buffer bytes.Buffer

	buffer.WriteString("curl ")
	if len(httpStream.Request.Host) > 1 {
		buffer.WriteString(httpStream.Request.Host)
	} else {
		buffer.WriteString(httpStream.ResponseEndpoint.Ip)
		buffer.WriteString(":")
		buffer.WriteString(fmt.Sprintf("%d", httpStream.ResponseEndpoint.Port))
	}

	buffer.WriteString(httpStream.Request.Uri)
	buffer.WriteString(" ")
	buffer.WriteString("-X")
	buffer.WriteString(httpStream.Request.Verb)
	buffer.WriteString(" ")
	for _, header := range httpStream.Request.Headers {
		buffer.WriteString(`-H "`)
		buffer.WriteString(header.Key)
		buffer.WriteString(": ")
		buffer.WriteString(header.Value)
		buffer.WriteString(`" `)
	}

	return buffer.String()
}
