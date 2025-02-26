package originmodel

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/light-pan/sharingan/replayer-agent/model/esmodel"
)

func RetrieveSessionId(data []byte) (esmodel.SessionId, error) {
	var source IDSource
	var id esmodel.SessionId
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, &source)
	if err != nil {
		return id, err
	}

	id = source.Data

	return id, nil
}

// 原始流量存储的sessionID数据格式
type IDSource struct {
	Data esmodel.SessionId `json:"data"`
}
