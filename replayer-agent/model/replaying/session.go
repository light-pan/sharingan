package replaying

import (
	"github.com/light-pan/sharingan/replayer-agent/model/recording"
)

type Session struct {
	Context         string
	SessionId       string
	CallFromInbound *recording.CallFromInbound
	ReturnInbound   *recording.ReturnInbound
	CallOutbounds   []*recording.CallOutbound
	RedirectDirs    map[string]string
	MockFiles       map[string][][]byte
	AppendFiles     []*recording.AppendFile
	TracePaths      []string
	ReadStorages    []*recording.ReadStorage
}

func NewSession() *Session {
	return &Session{}
}
