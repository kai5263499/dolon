package interfaces

import (
	"github.com/kai5263499/dolon/types"
)

type HttpProcessor interface {
	StartProcessingTcpSessions(Source)
	HttpSessionChan() chan *types.HttpSession
}
