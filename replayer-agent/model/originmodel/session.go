package originmodel

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/light-pan/sharingan/replayer-agent/model/esmodel"
)

func RetrieveSession(data []byte) (esmodel.Session, error) {
	var source DataSource
	var session esmodel.Session
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, &source)
	if err != nil {
		return session, err
	}
	session = source.Data
	return session, nil
}

// 原始流量的完整数据格式
type DataSource struct {
	Data esmodel.Session `json:"data"`
}
