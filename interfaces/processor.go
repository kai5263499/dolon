package interfaces

import (
	"github.com/kai5263499/dolon/types"
)

// HttpProcessor takes TCP sessions and returns
type HttpProcessor interface {
	StartProcessingTcpSessions(Source)
	HttpSessionChan() chan *types.HttpSession
}
