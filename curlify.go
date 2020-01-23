package dolon

import (
	"bytes"

	types "github.com/kai5263499/dolon/types"
)

// Curlify takes an HttpSession and returns a curl command that approxiamtes the http request
func Curlify(httpSession *types.HttpSession) string {
	var buffer bytes.Buffer

	buffer.WriteString("curl ")
	if len(httpSession.Request.Host) > 1 {
		buffer.WriteString(httpSession.Request.Host)
	}

	buffer.WriteString(httpSession.Request.Uri)
	buffer.WriteString(" ")
	buffer.WriteString("-X")
	buffer.WriteString(httpSession.Request.Verb)
	buffer.WriteString(" ")
	for _, header := range httpSession.Request.Headers {
		buffer.WriteString(`-H "`)
		buffer.WriteString(header.Key)
		buffer.WriteString(": ")
		buffer.WriteString(header.Value)
		buffer.WriteString(`" `)
	}

	return buffer.String()
}
