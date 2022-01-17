/*
config.jsonを読み込み、中身を保持する機能を提供するパッケージ
*/

package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"webapi/utils"
)

var serversConfigPath string

func SetConfPath(confName string) {
	currentDir, err := utils.GetCurrentDir()
	if err != nil {
		panic("err msg: " + err.Error())
	}
	serversConfigPath = filepath.Join(currentDir, confName)
}

type serversConfig struct {
	APIGateWayServers []string `json:"APIGateWayServers"`
}

// NewServerConfig はservers.jsonの中身をserversConfig構造体にセットし、返す
func NewServerConfig() *serversConfig {
	if serversConfigPath == "" {
		SetConfPath("config.json")
	}
	// 構造体を初期化
	conf := &serversConfig{}

	// 設定ファイルを読み込む
	cValue, err := ioutil.ReadFile(serversConfigPath)
	if err != nil {
		panic(err.Error())
	}

	// 読み込んだjson文字列をデコードし構造体にマッピング
	err = json.Unmarshal([]byte(cValue), conf)
	if err != nil {
		panic(err.Error())
	}

	return conf
}
